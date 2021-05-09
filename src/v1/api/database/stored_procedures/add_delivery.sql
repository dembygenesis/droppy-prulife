/**
  Transactions
 */

DROP PROCEDURE IF EXISTS  `add_delivery`;
DELIMITER $$

CREATE PROCEDURE `add_delivery` (
    p_delivery_option TEXT,
    p_seller_id INTEGER,
    p_dropshipper_id INTEGER,
    p_name TEXT,
    p_contact_number TEXT,
    p_address TEXT,
    p_note TEXT,
    p_region TEXT,
    p_service_fee DECIMAL,
    p_declared_amount DECIMAL,
    p_delivery_details TEXT
)
BEGIN

    DECLARE p_base_price INTEGER;

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;  -- rollback any changes made in the transaction
            RESIGNAL;  -- raise again the sql exception to the caller
        END;



    START TRANSACTION;

    -- Validate declared amount
    IF p_declared_amount <= 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Declared amount must be greater than 0';
    END IF;

    -- Modify declared amount to then be adjusted to service fee
    SET p_service_fee = (SELECT get_service_fee(p_declared_amount));
    SET p_base_price = (SELECT get_base_price(p_declared_amount));

    -- Validate additional details
    IF p_name = '' OR p_contact_number = '' OR p_note = '' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Buyer informations: name, contact_number, or note must not be empty.';
    END IF;

    -- Validate p_delivery_option
    SELECT
        COUNT(id),
        id,
        `name`
    INTO @delivery_option_exists, @delivery_option_id, @delivery_option
    FROM `delivery_option`
    WHERE 1 = 1
      AND name = p_delivery_option;

    IF @delivery_option_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Delivery type does not exist!';
    END IF;

    -- Validate region.
    SET @region_exists = (
        SELECT
            COUNT(id)
        FROM
            `region`
        WHERE 1 = 1
          AND `name` = p_region
    );

    IF @region_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Region does not exist';
    END IF;

    SELECT
        `id`,
        `name`
    INTO @region_id, @region_name
    FROM region
    WHERE 1 = 1
      AND `name` = p_region;

    /**
      Use system based config and override the dropshipper ID to use the system based rules
     */

    -- SELECT @region_name, @delivery_option;

    -- SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Test Break Here';

    IF @region_name = 'Luzon' AND @delivery_option = 'Parcel' THEN
        SET p_dropshipper_id = (
            SELECT
                id
            FROM `user` u
            WHERE 1 = 1
              AND email = (SELECT `value` FROM sysparam WHERE `key` = 'HANDLER_PARCEL_LUZON')
        );
    END IF;

    IF @region_name = 'Luzon' AND @delivery_option = 'Dropship' THEN
        SET p_dropshipper_id = (
            SELECT
                id
            FROM `user` u
            WHERE 1 = 1
              AND email = (SELECT `value` FROM sysparam WHERE `key` = 'HANDLER_DROPSHIP_LUZON')
        );
    END IF;

    IF @region_name = 'Vis/Min' AND @delivery_option = 'Parcel' THEN
        SET p_dropshipper_id = (
            SELECT
                id
            FROM `user` u
            WHERE 1 = 1
              AND email = (SELECT `value` FROM sysparam WHERE `key` = 'HANDLER_PARCEL_VISMIN')
        );
    END IF;

    IF @region_name = 'Vis/Min' AND @delivery_option = 'Dropship' THEN
        SET p_dropshipper_id = (
            SELECT
                id
            FROM `user` u
            WHERE 1 = 1
              AND email = (SELECT `value` FROM sysparam WHERE `key` = 'HANDLER_DROPSHIP_VISMIN')
        );
    END IF;

    -- Validate seller.
    /*SET @user_exists = (
        SELECT COUNT(id) FROM `user`
        WHERE 1 = 1
          AND is_active = 1
          AND id = p_seller_id
          AND user_type_id = (SELECT `id` FROM user_type WHERE `name` = 'Seller')
    );*/

    -- New seller query
    SELECT
        COUNT(u.id),
        ut.coin_amount
    INTO @user_exists, @user_coin_balance
    FROM `user` u
             INNER JOIN user_total ut
                        ON  1 = 1
                            AND u.id = ut.user_id
    WHERE 1 = 1
      AND u.is_active = 1
      AND u.id = p_seller_id
      AND u.user_type_id = (SELECT `id` FROM user_type WHERE `name` = 'Seller');

    IF @user_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Seller does not exist';
    END IF;

    -- Validate dropshipper.
    SET @dropshipper_exists = (
        SELECT COUNT(id) FROM `user`
        WHERE 1 = 1
          AND is_active = 1
          AND id = p_dropshipper_id
          AND user_type_id = (SELECT `id` FROM user_type WHERE `name` = 'Dropshipper')
    );

    IF @dropshipper_exists = 0 THEN
        SET @error_message = CONCAT('Dropshipper does not exist for user id: ', @dropshipper_exists, ' charles ', p_dropshipper_id);

        -- SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Dropshipper does not exist';
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = @error_message;
    END IF;

    -- Validate that the seller has ENOUGH coins that he'd be able to make that transaction
    IF @user_coin_balance < p_service_fee THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'You dont have enough coins.';
    END IF;



    -- Validate delivery details
    -- Loop orders and get total lol -- add price column per quantity
    SET @str_length = LENGTH(`p_delivery_details`);
    SET @iterator = 1;
    SET @comma_count = 1;
    SET @last_comma_pos = 1;
    SET @value = '';

    -- Create a temp table to house your orders variables
    SET @unique_tmp_table_name = (
        SELECT CONCAT("table_prefix_", REPLACE(UUID(), '-', '')) AS haha
    );

    SET @unique_tmp_table_name_create_stmt = CONCAT('
        CREATE TEMPORARY TABLE IF NOT EXISTS ', @unique_tmp_table_name, ' (
          product_id INTEGER PRIMARY KEY,
          quantity INTEGER,
          price_per_item_distributor DECIMAL(65,2),
          total_price_distributor DECIMAL(65,2)
        ) ;
    ');

    PREPARE stmt FROM @unique_tmp_table_name_create_stmt;
    EXECUTE stmt;

    /**
      Process products via loop and insert them to a temporary table to so we can manipulate them better
     */

    product_validation:
    LOOP
        IF @iterator > @str_length  THEN
            LEAVE product_validation;
        END IF;

        IF SUBSTRING(`p_delivery_details`, @iterator, 1) = ',' OR SUBSTRING(`p_delivery_details`, @iterator + 1, 1) = '' THEN

            -- Base Case: First comma found
            IF (@comma_count = 1) THEN
                SET @value = SUBSTRING_INDEX(`p_delivery_details`, ',', @comma_count);

                SET @comma_count = @comma_count + 1;
                SET @last_comma_pos = @iterator;
            ELSE
                -- Substring index here
                SET @value = SUBSTRING_INDEX(`p_delivery_details`, ',', @comma_count);
                SET @value = SUBSTRING(@value, @last_comma_pos + 1, @iterator);

                SET @comma_count = @comma_count + 1;
                SET @last_comma_pos = @iterator;
            END IF;

            -- Fetch update variables
            SET @product_id = SUBSTRING_INDEX(@value, '---', 1);
            SET @quantity = SUBSTRING(
                    SUBSTRING_INDEX(@value, '---', 2),
                    LOCATE('---', @value) + 3,
                    LENGTH(@value)
                );

            -- Validate product if existing (active)
            SET @validate_product_stmt = CONCAT('
              SELECT
                COUNT(id),
                name,
                price_per_item_dropshipper
              INTO @product_count, @product_name, @product_price_per_item_dropshipper
              FROM product
              WHERE 1 = 1
                AND id = ', @product_id, '
                AND is_active = 1
            ');

            PREPARE stmt FROM @validate_product_stmt;
            EXECUTE stmt;

            IF @product_count = 0 THEN
                SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'A product does not exist';
            END IF;

            /**
              Droppy custom logic
             */
            -- Disallow max cee blister with dropship
            IF p_delivery_option = 'Dropship' AND @product_name = 'Max-Cee Blister' THEN
                SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Max-Cee Blister is not available for Dropship.';
            END IF;

            IF p_delivery_option = 'Parcel' AND @product_name = 'Max-Cee' THEN
                SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Max-Cee is not available for Parcel.';
            END IF;

            -- Validate quantity in inventory if parcel
            IF p_delivery_option = 'Parcel' THEN
                SET @has_enough_products = (
                    SELECT
                        COUNT(id)
                    FROM inventory i
                    WHERE 1 = 1
                      AND i.product_id = @product_id
                      AND i.quantity >= @quantity
                      AND i.region_id = @region_id
                      AND i.seller_id = p_seller_id
                      AND i.dropshipper_id = p_dropshipper_id
                );

                /*SELECT
                       @product_id,
                       @quantity,
                       @region_id,
                       p_seller_id,
                       p_dropshipper_id;*/

                IF @has_enough_products = 0 THEN
                    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Product quantity invalid. You do not have enough product to make this delivery.';
                END IF;
            END IF;

            -- Insert products into temporary table
            SET @insert_tmp_table_stmt = CONCAT('
              INSERT INTO ', @unique_tmp_table_name, '
              SELECT ', @product_id,', ', @quantity, ', ', @product_price_per_item_dropshipper, ', ', @quantity * @product_price_per_item_dropshipper, '
            ');

            PREPARE stmt FROM @insert_tmp_table_stmt;
            EXECUTE stmt;

        END IF;

        SET @iterator = @iterator + 1;

    END LOOP product_validation;

    -- Extract variables here
    SET @extract_totals = CONCAT('
        SELECT
          SUM(total_price_distributor)
        INTO @amount_distributor
        FROM ', @unique_tmp_table_name, '
    ');

    PREPARE stmt FROM @extract_totals;
    EXECUTE stmt;

    -- IF PARCEL: Sync with inventory (if dropship just proceed)
    IF p_delivery_option = 'Parcel' THEN
        -- Update inventory seller
        SET @update_inventory_stmt = CONCAT('
            UPDATE inventory i
            INNER JOIN ', @unique_tmp_table_name, ' a
                ON i.product_id = a.product_id
            SET i.quantity = i.quantity - a.quantity
            WHERE 1 = 1
                AND i.region_id = ', @region_id, '
                AND i.seller_id = ', p_seller_id, '
                AND i.dropshipper_id = ', p_dropshipper_id, '
        ');

        PREPARE stmt FROM @update_inventory_stmt;
        EXECUTE stmt;

        -- Update inventory dropshipper
    END IF;

    -- Ensure there are enough coins on the user to make this order


    -- INSERT INTO details
    INSERT INTO `delivery` (created_by,
                            created_date,
                            is_active,
                            `name`,
                            `address`,
                            region_id,
                            service_fee,
                            base_price,
                            declared_amount,
                            delivery_option_id,
                            seller_id,
                            dropshipper_id,
                            delivery_status_id,
                            contact_number,
                            note,
                            amount_distributor
    )
    VALUES (
               p_seller_id,
               NOW(),
               1,
               p_name,
               p_address,
               @region_id,
               p_service_fee,
               p_base_price,
               p_declared_amount,
               @delivery_option_id,
               p_seller_id,
               p_dropshipper_id,
               (SELECT id FROM delivery_status WHERE name = 'Proposed'),
               p_contact_number,
               p_note,
               @amount_distributor
           );

    SET @delivery_id_created = LAST_INSERT_ID();

    -- We use a dynamic insert to get the dynamic table name created
    SET @update_delivery_detail_stmt = CONCAT('
        INSERT INTO `delivery_detail` (delivery_id, product_id, quantity, price_per_item_distributor, total_price_distributor)
        (SELECT
            ', @delivery_id_created, ',
            product_id,
            quantity,
            price_per_item_distributor,
            total_price_distributor
        FROM ', @unique_tmp_table_name, ')
    ');

    PREPARE stmt FROM @update_delivery_detail_stmt;
    EXECUTE stmt;

    -- Add to delivery tracking
    INSERT INTO delivery_tracking (delivery_id, delivery_status_id, last_updated, updated_by)
    VALUES (
               @delivery_id_created,
               (SELECT id FROM delivery_status WHERE name = 'Proposed'),
               NOW(),
               p_seller_id
           );

    -- Update user total to also deduct service fee
    UPDATE user_total
    SET coin_amount = coin_amount - p_service_fee
    WHERE user_id = p_seller_id;

    INSERT INTO coin_transaction (created_by, created_date, is_active, user_id, type, amount, delivery_id)
    VALUES (
               @admin_account,
               NOW(),
               1,
               p_seller_id,
               'D',
               ABS(p_service_fee),
               @delivery_id_created
           );

    SET @admin_account = (
        SELECT MAX(id)
        FROM `user`
        WHERE 1 = 1
          AND is_active = 1
          AND user_type_id = (SELECT `id` FROM user_type WHERE `name` = 'Admin')
    );


    /**
      Profits for dropshipper and admin
     */
    -- Add 35 coins to dropshipper (p_service_fee - 35)
    SET @dropshipper_totals_exists = (
        SELECT COUNT(u.id) FROM `user` u
                                    INNER JOIN user_total ut
                                               ON 1 = 1
                                                   AND u.id = ut.user_id
        WHERE 1 = 1
          AND u.is_active = 1
          AND u.id = p_dropshipper_id
          AND u.user_type_id = (SELECT `id` FROM user_type WHERE `name` = 'Dropshipper')
    );

    IF @dropshipper_totals_exists = 0 THEN
        INSERT INTO `user_total` (user_id, amount, coin_amount, created_by, last_updated)
        VALUES (
                   p_dropshipper_id,
                   0,
                   35,
                   @admin_account,
                   NOW()
               );
    ELSE
        UPDATE user_total
        SET coin_amount = coin_amount + 35
        WHERE user_id = p_dropshipper_id;
    END IF;

    INSERT INTO coin_transaction (created_by, created_date, is_active, user_id, type, amount, delivery_id)
    VALUES (
               @admin_account,
               NOW(),
               1,
               p_dropshipper_id,
               'C',
               -35,
               @delivery_id_created
           );

    -- Add (p_service_fee - 35) coins to admin
    SET @admin_totals_exists = (
        SELECT COUNT(u.id) FROM `user` u
                                    INNER JOIN user_total ut
                                               ON 1 = 1
                                                   AND u.id = ut.user_id
        WHERE 1 = 1
          AND u.is_active = 1
          AND u.id = @admin_account
          AND u.user_type_id = (SELECT `id` FROM user_type WHERE `name` = 'Admin')
    );

    IF @admin_totals_exists = 0 THEN
        INSERT INTO `user_total` (user_id, amount, coin_amount, created_by, last_updated)
        VALUES (
                   @admin_account,
                   0,
                   (p_service_fee - 35),
                   @admin_account,
                   NOW()
               );
    ELSE
        UPDATE user_total
        SET coin_amount = coin_amount + (p_service_fee - 35)
        WHERE user_id = @admin_account;
    END IF;

    INSERT INTO coin_transaction (created_by, created_date, is_active, user_id, type, amount, delivery_id)
    VALUES (
               @admin_account,
               NOW(),
               1,
               @admin_account,
               'D',
               (p_service_fee - 35),
               @delivery_id_created
           );

    -- SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Set delivery test';
    COMMIT;
END $$

DELIMITER;

