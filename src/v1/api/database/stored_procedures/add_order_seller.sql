/**
  Transactions
 */

DROP PROCEDURE IF EXISTS  `add_order_seller`;
DELIMITER $$

CREATE PROCEDURE `add_order_seller` (
    -- person authentication first
    -- @p_admin_id INTEGER,
    p_user_id INTEGER,
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

    -- Validate Seller
    SET @user_exists = (
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
          AND user.id = p_user_id
          AND user.is_active = 1
    );

    IF @user_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Seller does not exist / or has not been loaded with any transaction';
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

    SET @region_id = (
        SELECT
            `id`
        FROM region
        WHERE 1 = 1
          AND `name` = p_region
    );

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
        AND user_id = p_user_id;

    IF @total_price > @total_coins_available THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'User cannot afford the products!';
    END IF;

    /**
      Process the transactions
      1. Create order
      2. Create order items
      3. Subtract from user totals
      4. Return the coins to droppy
     */

    -- Create Order
    INSERT INTO `order` (created_by, user_id, created_date, is_active, order_status_id, amount, region_id)
    VALUES (
        p_user_id,
        p_user_id,
        NOW(),
        1,
        (SELECT id FROM order_status WHERE `name` = 'Proposed'),
        @total_price,
        @region_id
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
            p_user_id,
            NOW(),
            1,
            @order_id_created,
            @product_id,
            @product_quantity,
            @product_price_per_item,
            @product_quantity * @product_price_per_item
        );

        -- Update inventory
        SET @inventory_exists_in_user = (
            SELECT
                COUNT(*)
            FROM inventory
            WHERE 1 = 1
                AND user_id = p_user_id
                AND product_id = @product_id
                AND region_id = @region_id
        );

        IF @inventory_exists_in_user > 0 THEN
            UPDATE inventory
                SET quantity = quantity + @product_quantity
            WHERE user_id = p_user_id
                AND product_id = @product_id
                AND region_id = @region_id;
        ELSE
            INSERT INTO inventory (product_id, quantity, created_by, user_id, created_date, is_active, region_id)
            VALUES (
                @product_id,
                @product_quantity,
                p_user_id,
                p_user_id,
                NOW(),
                1,
                @region_id
            );
        END IF;

        SET @iterator = @iterator + 1;
    END WHILE;

    -- Subtract from user totals
    UPDATE user_total
        SET coin_amount = coin_amount - @total_price
    WHERE user_id = p_user_id;

    -- Return the coins to droppy (use admin ID) - use default for this in the backend
    INSERT INTO coin_transaction (created_by, created_date, is_active, order_id, `type`, amount)
    VALUES (
        @p_admin_id,
        NOW(),
        1,
        @order_id_created,
        'D',
        @total_price
    );

    UPDATE user_total
        SET coin_amount = coin_amount + @total_price
    WHERE user_id = @p_admin_id;

    COMMIT;
END $$

DELIMITER;

