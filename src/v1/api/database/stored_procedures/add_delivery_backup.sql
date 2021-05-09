/**
  Transactions
 */

DROP PROCEDURE IF EXISTS  `add_delivery`;
DELIMITER $$

CREATE PROCEDURE `add_delivery` (
    p_delivery_option TEXT,
    p_user_id INTEGER,
    p_name TEXT,
    p_address TEXT,
    p_region TEXT,
    p_service_fee DECIMAL,
    p_declared_amount DECIMAL,
    p_delivery_details TEXT
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;  -- rollback any changes made in the transaction
            RESIGNAL;  -- raise again the sql exception to the caller
        END;

    START TRANSACTION;

    -- Validate p_delivery_option
    SELECT
        COUNT(id),
        id
    INTO @delivery_option_exists, @delivery_option_id
    FROM `delivery_option`
    WHERE 1 = 1
      AND name = p_delivery_option;

    IF @delivery_option_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Delivery type does not exist!';
    END IF;

    -- Validate user.
    SET @user_exists = (
        SELECT MAX(id) FROM `user`
        WHERE 1 = 1
          AND is_active = 1
          AND id = p_user_id
          AND user_type_id = (SELECT `id` FROM user_type WHERE `name` = 'Seller')
    );

    IF @user_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Seller does not exist';
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

    SET @region_id = (
        SELECT
            id
        FROM
            `region`
        WHERE 1 = 1
          AND `name` = p_region
    );

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
                COUNT(id) INTO @product_count
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
                        AND i.user_id = p_user_id
                );

                IF @has_enough_products = 0 THEN
                    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Product quantity invalid. You do not have enough product to make this delivery.';
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

    PREPARE stmt FROM @total_price_stmt;
    EXECUTE stmt;

    -- IF PARCEL: Sync with inventory (if dropship just proceed)
    IF p_delivery_option = 'Parcel' THEN
        -- Get total price
        SET @update_inventory_stmt = CONCAT('
            UPDATE inventory i
            INNER JOIN ', @unique_tmp_table_name, ' a
                ON i.product_id = a.product_id
            SET i.quantity = i.quantity - a.quantity
            WHERE 1 = 1
                AND i.region_id = ', @region_id, '
                AND i.user_id = ', p_user_id, '
        ');

        PREPARE stmt FROM @update_inventory_stmt;
        EXECUTE stmt;
    END IF;

    -- INSERT INTO details
    INSERT INTO `delivery` (created_by,
                            created_date,
                            is_active,
                            `name`,
                            `address`,
                            region_id,
                            service_fee,
                            declared_amount,
                            delivery_option_id
                            )
    VALUES (
        p_user_id,
        NOW(),
        1,
        p_name,
        p_address,
        @region_id,
        p_service_fee,
        p_declared_amount,
        @delivery_option_id
    );

    SET @delivery_id_created = LAST_INSERT_ID();

    -- We use a dynamic insert to get the dynamic table name created
    SET @update_delivery_detail_stmt = CONCAT('
        INSERT INTO `delivery_detail` (delivery_id, product_id, quantity)
        (SELECT
            ', @delivery_id_created, ',
            product_id,
            quantity
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
        p_user_id
    );

    COMMIT;
END $$

DELIMITER;

