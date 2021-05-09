DROP PROCEDURE IF EXISTS  `add_order`;
DELIMITER $$

CREATE PROCEDURE `add_order` (
    p_seller_id INTEGER,
    p_dropshipper_id INTEGER,
    p_order_details TEXT,
    p_region TEXT
)
BEGIN

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;  -- rollback any changes made in the transaction
            RESIGNAL;  -- raise again the sql exception to the caller
        END;

    START TRANSACTION;

    /**
      Note!!!
      This is just grabbing the latest admin created. (lol)
     */

    -- Validate default admin.
    SET @p_admin_id = (
        SELECT MAX(id) FROM `user`
        WHERE 1 = 1
          AND is_active = 1
          AND user_type_id = (SELECT `id` FROM user_type WHERE `name` = 'Admin')
    );

    IF @admin_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Admin does not exist';
    END IF;

    -- Validate Region
    SET @region_exists = (
        SELECT
            COUNT(id)
        FROM region
        WHERE 1 = 1
          AND `name` = p_region
    );

    IF @region_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Region does not exist!';
    END IF;

    -- Set region_id and region_name into variables
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
    IF @region_name = 'Luzon' THEN
        SET p_dropshipper_id = (
            SELECT
                id
            FROM `user` u
            WHERE 1 = 1
                AND email = (SELECT `value` FROM sysparam WHERE `key` = 'HANDLER_PACKAGE_LUZON')
        );
    END IF;

    IF @region_name = 'Vis/Min' THEN
        SET p_dropshipper_id = (
            SELECT
                id
            FROM `user` u
            WHERE 1 = 1
              AND email = (SELECT `value` FROM sysparam WHERE `key` = 'HANDLER_PACKAGE_VISMIN')
        );
    END IF;

    -- Validate Seller
    SET @seller_exists = (
        SELECT
            COUNT(user.id)
        FROM user
        INNER JOIN user_type
           ON 1 = 1
               AND user.user_type_id = user_type.id
        INNER JOIN user_total
           ON 1 = 1
               AND user_total.user_id = user.id
        WHERE 1 = 1
          AND is_active = 1
          AND user_type.name = 'Seller'
          AND user.id = p_seller_id
          AND user.is_active = 1
    );

    IF @seller_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Seller does not exist / or has not been loaded with any transaction';
    END IF;

    -- Validate Dropshipper
    SET @dropshipper_exists = (
        SELECT
            COUNT(user.id)
        FROM user
                 INNER JOIN user_type
                            ON 1 = 1
                                AND user.user_type_id = user_type.id
        WHERE 1 = 1
          AND is_active = 1
          AND user_type.name = 'Dropshipper'
          AND user.id = p_dropshipper_id
          AND user.is_active = 1
    );

    IF @dropshipper_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Dropshipper does not exist / or has not been loaded with any transaction';
    END IF;

    -- Loop orders and get total lol -- add price column per quantity
    SET @str_length = LENGTH(`p_order_details`);
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
          quantity INTEGER
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

        IF SUBSTRING(`p_order_details`, @iterator, 1) = ',' OR SUBSTRING(`p_order_details`, @iterator + 1, 1) = '' THEN

            -- Base Case: First comma found
            IF (@comma_count = 1) THEN
                SET @value = SUBSTRING_INDEX(`p_order_details`, ',', @comma_count);

                SET @comma_count = @comma_count + 1;
                SET @last_comma_pos = @iterator;
            ELSE
                -- Substring index here
                SET @value = SUBSTRING_INDEX(`p_order_details`, ',', @comma_count);
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
                name
              INTO @product_count, @product_name
              FROM product
              WHERE 1 = 1
                AND id = ', @product_id, '
                AND is_active = 1
            ');

            PREPARE stmt FROM @validate_product_stmt;
            EXECUTE stmt;

            -- Validate product exists
            IF @product_count = 0 THEN
                SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'A product does not exist';
            END IF;

            /**
              Droppy custom rules
             */
            IF @product_name = 'Max-Cee' THEN
                SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Max-Cee is only available for Dropship.';
            END IF;

            IF @product_name = 'Max-Cee Blister' THEN
                IF @quantity < 10 THEN
                    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'You can only order Max-Cee Blisters 10 at a time.';
                END IF;

                IF MOD(@quantity, 10) != 0 THEN
                    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'You can only order Max-Cee Blisters 10 at a time.';
                END IF;
            END IF;

            -- Insert products into temporary table
            SET @insert_tmp_table_stmt = CONCAT('
              INSERT INTO ', @unique_tmp_table_name, '
              SELECT ', @product_id,', ', @quantity, '
            ');

            PREPARE stmt FROM @insert_tmp_table_stmt;
            EXECUTE stmt;

        END IF;

        SET @iterator = @iterator + 1;

    END LOOP product_validation;


    SET @check_inserted_products_stmt = CONCAT('
        SELECT COUNT(*) INTO @product_check FROM ', @unique_tmp_table_name, '
    ');

    PREPARE stmt FROM @check_inserted_products_stmt;
    EXECUTE stmt;

    IF @product_check = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'No products were detected, please check your details.';
    END IF;

    /**
      Do further validations here
     */

    -- Get total price
    SET @total_price_stmt = CONCAT('
      SELECT
        SUM(p.price_per_item * a.quantity),
        GROUP_CONCAT(CONCAT(
            p.id,
            \'-\',
            a.quantity,
            \'-\',
            p.price_per_item)
        )
        INTO @total_price, @product_ids
      FROM ', @unique_tmp_table_name, ' a
      INNER JOIN product p
        ON 1 = 1
            AND a.product_id = p.id
    ');

    PREPARE stmt FROM @total_price_stmt;
    EXECUTE stmt;

    -- Compare against user coins
    SELECT
        coin_amount INTO @total_coins_available
    FROM user_total
    WHERE 1 = 1
      AND user_id = p_seller_id;

    IF @total_price > @total_coins_available THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'User cannot afford the products!';
    END IF;

    /**
      M88 Logic
      1. Must not exceed PACKAGE_PURCHASE_LIMIT var (currently 10750 as of aug 10, 2020)
      2. Must always be deducted by PACKAGE_PURCHASE_DEDUCTION var per transactions (currently 10750 as of aug 10, 2020)
     */

    SET @PACKAGE_PURCHASE_DEDUCTION = (SELECT CAST(value AS DECIMAL(65,2)) FROM sysparam WHERE `key` = 'PACKAGE_PURCHASE_DEDUCTION');
    SET @PACKAGE_PURCHASE_LIMIT = (SELECT CAST(value AS DECIMAL(65,2)) FROM sysparam WHERE `key` = 'PACKAGE_PURCHASE_LIMIT');

    -- Must not exceed "PACKAGE_PURCHASE_LIMIT"
    IF @total_price > @PACKAGE_PURCHASE_LIMIT THEN
        SET @err_message = CONCAT('You can only purchase products up to a max worth of: ', @PACKAGE_PURCHASE_LIMIT);

        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = @err_message;
    END IF;

    IF @total_coins_available < @PACKAGE_PURCHASE_LIMIT THEN
        SET @err_message = CONCAT('You must have at least : ', @PACKAGE_PURCHASE_DEDUCTION, ' in order to make a purchase.');

        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = @err_message;
    END IF;

     /**
       end of M88 Logic
      */

    -- Create Order
    INSERT INTO `order` (created_by, seller_id, dropshipper_id, created_date, is_active, order_status_id, amount, region_id, admin_id)
    VALUES (
               p_seller_id,
               p_seller_id,
               p_dropshipper_id,
               NOW(),
               1,
               (SELECT id FROM order_status WHERE `name` = 'Proposed'),
               -- The actual price replaced by the minimum purchase price
               -- @total_price,
               @PACKAGE_PURCHASE_DEDUCTION,
               @region_id,
               @p_admin_id
           );

    SET @order_id_created = LAST_INSERT_ID();

    -- Create order items
    SET @iterator = 1;
    SET @current_data = '';

    SET @product_order_count = LENGTH(@product_ids) - LENGTH(REPLACE(@product_ids, ',', '')) + 1;

    WHILE (@iterator <= @product_order_count) DO
            SELECT
                SUBSTRING_INDEX( SUBSTRING_INDEX(@product_ids, ',', @iterator ), ',', -1 )
            INTO @product_id_quantity_price;

            SET @product_id = (
                SELECT
                    SUBSTRING_INDEX(
                            SUBSTRING_INDEX(@product_id_quantity_price, '-', 1),
                            '-',
                            -1
                        )
            );

            SET @product_quantity = (
                SELECT
                    SUBSTRING_INDEX(
                            SUBSTRING_INDEX(@product_id_quantity_price, '-', 2),
                            '-',
                            -1
                        )
            );

            SET @product_price_per_item = (
                SELECT
                    SUBSTRING_INDEX(
                            SUBSTRING_INDEX(@product_id_quantity_price, '-', 3),
                            '-',
                            -1
                        )
            );

            -- Add to order detail.
            INSERT INTO `order_detail` (
                `created_by`,
                `created_date`,
                `is_active`,
                `order_id`,
                `product_id`,
                `quantity`,
                `price_per_item`,
                `total_price`
            )
            VALUES
            (
                p_seller_id,
                NOW(),
                1,
                @order_id_created,
                @product_id,
                @product_quantity,
                @product_price_per_item,
                @product_quantity * @product_price_per_item
            );

            -- Update inventory for both seller and dropshipper
            SET @inventory_exists_in_seller = (
                SELECT
                    COUNT(*)
                FROM inventory
                WHERE 1 = 1
                  AND seller_id = p_seller_id
                  AND dropshipper_id = p_dropshipper_id
                  AND product_id = @product_id
                  AND region_id = @region_id
            );

            IF @inventory_exists_in_seller > 0 THEN
                UPDATE inventory
                SET quantity = quantity + @product_quantity
                WHERE seller_id = p_seller_id
                  AND dropshipper_id = p_dropshipper_id
                  AND product_id = @product_id
                  AND region_id = @region_id;
            ELSE
                INSERT INTO inventory (product_id, quantity, created_by, seller_id, dropshipper_id, created_date, is_active, region_id)
                VALUES (
                           @product_id,
                           @product_quantity,
                           p_seller_id,
                           p_seller_id,
                           p_dropshipper_id,
                           NOW(),
                           1,
                           @region_id
                       );
            END IF;

            SET @iterator = @iterator + 1;
        END WHILE;

    /**
      SELLER (subtract full amount)
     */

    INSERT INTO coin_transaction (created_by, created_date, user_id, is_active, order_id, `type`, amount)
    VALUES (
               @p_admin_id,
               NOW(),
               p_seller_id,
               1,
               @order_id_created,
               'D',
               ABS(@PACKAGE_PURCHASE_DEDUCTION)
           );

    UPDATE user_total
    SET coin_amount = amount - ABS(@PACKAGE_PURCHASE_DEDUCTION)
    WHERE user_id = p_seller_id;

    /**
      DROPSHIPPER (subtract full amount)
     */

    -- Add 75 coins to dropshipper and update his totals
    INSERT INTO coin_transaction (created_by, created_date, user_id, is_active, order_id, `type`, amount)
    VALUES (
               @p_admin_id,
               NOW(),
               p_dropshipper_id,
               1,
               @order_id_created,
               'C',
               -75
           );

    UPDATE user_total
    SET coin_amount = coin_amount + 75
    WHERE user_id = p_dropshipper_id;

    /**
      ADMIN (full amount - 75)
     */

    -- Add (full amount - 75) to admin
    INSERT INTO coin_transaction (created_by, created_date, user_id, is_active, order_id, `type`, amount)
    VALUES (
               @p_admin_id,
               NOW(),
               @p_admin_id,
               1,
               @order_id_created,
               'D',
               ABS(@PACKAGE_PURCHASE_DEDUCTION) - 75
           );

    UPDATE user_total
    SET coin_amount = coin_amount + ABS(@PACKAGE_PURCHASE_DEDUCTION) - 75
    WHERE user_id = @p_admin_id;

    COMMIT;
END $$

DELIMITER;

