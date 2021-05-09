DROP PROCEDURE IF EXISTS  `get_service_fee`;
DELIMITER $$

CREATE PROCEDURE `get_service_fee` (
    p_order_details TEXT
)
this_proc:BEGIN
    DECLARE service_fee DECIMAL;
    DECLARE base_price DECIMAL;

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;  -- rollback any changes made in the transaction
            RESIGNAL;  -- raise again the sql exception to the caller
        END;

    START TRANSACTION;

    IF p_order_details = "" OR p_order_details = '""' THEN
        SET service_fee = 0;
        SELECT service_fee;
        LEAVE this_proc;
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
                COUNT(id) INTO @product_count
              FROM product
              WHERE 1 = 1
                AND id = ', @product_id, '
                AND is_active = 1
            ');

            PREPARE stmt FROM @validate_product_stmt;
            EXECUTE stmt;

            IF @product_count = 0 THEN
                SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'You have no valid products.';
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

    IF @total_price BETWEEN 0 AND 1499 THEN
        SET service_fee = 195;
    ELSEIF @total_price BETWEEN 1500 AND 1999 THEN
        SET service_fee = 200;
    ELSEIF @total_price BETWEEN 2000 AND 2499 THEN
        SET service_fee = 205;
    ELSEIF @total_price BETWEEN 2500 AND 2999 THEN
        SET service_fee = 240;
    ELSEIF @total_price BETWEEN 3000 AND 3499 THEN
        SET service_fee = 245;
    ELSEIF @total_price BETWEEN 3500 AND 3999 THEN
        SET service_fee = 250;
    ELSEIF @total_price BETWEEN 4000 AND 4499 THEN
        SET service_fee = 255;
    ELSEIF (@total_price BETWEEN 4500 AND 4999) OR @total_price > 4999 THEN
        SET service_fee = 260;
    END IF;

    SELECT service_fee;

    COMMIT;
END $$

DELIMITER;

