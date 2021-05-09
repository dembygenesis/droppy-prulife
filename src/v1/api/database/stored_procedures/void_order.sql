/**
  Transactions
 */

DROP PROCEDURE IF EXISTS  `void_order`;
DELIMITER $$

CREATE PROCEDURE `void_order` (
    -- person authentication first
    p_admin_id INTEGER,
    p_order_id INTEGER
)
BEGIN

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;  -- rollback any changes made in the transaction
        RESIGNAL;  -- raise again the sql exception to the caller
    END;

    START TRANSACTION;

    -- Validate admin
    SET @admin_exists = (
        SELECT
            COUNT(user.id)
        FROM user
        INNER JOIN user_type
            ON 1 = 1
                AND user.user_type_id = user_type.id
        WHERE 1 = 1
          AND is_active = 1
          AND user_type.name = 'Admin'
          AND user.id = p_admin_id
          AND user.is_active = 1
    );

    IF @admin_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Admin does not exist';
    END IF;

    -- Validate order
    SELECT
        COUNT(id),
        admin_id,
        seller_id,
        dropshipper_id,
        amount,
        region_id
    INTO @order_id_exists, @admin_id, @seller_id, @dropshipper_id, @amount, @region_id
    FROM `order`
    WHERE 1 = 1
      AND id = p_order_id
      AND is_active = 1;

    IF @order_id_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Order ID does not exist';
    END IF;

    /**
      Process the transactions
      1. Void order
      2. Void order items (automatic when voiding orders)
      3. Add coins back to user totals
      4. Subtract the coins from droppy (track this by created_by)
     */

    -- Void order
    UPDATE `order`
        SET is_active = 0,
            updated_by = p_admin_id,
            last_updated = NOW()
    WHERE id = p_order_id
        AND id = p_order_id;

    -- Update seller coin totals
    UPDATE user_total
        SET coin_amount = coin_amount + ABS(@amount)
    WHERE user_id = @seller_id
        AND last_updated = NOW()
        AND updated_by = p_admin_id;

    -- Update admin coin totals
    UPDATE user_total
    SET coin_amount = coin_amount - ABS(@amount)
    WHERE user_id = @admin_id;

    -- I should update dropshipper coin totals lol (TO DO)

    /**
      Just update 1 INVENTORY
     */

    -- Update seller inventory
    UPDATE
        inventory i
            INNER JOIN order_detail od
            ON i.`product_id` = od.`product_id`
            INNER JOIN `order` o
            ON od.order_id = o.id
    SET i.quantity = i.quantity - od.`quantity`
    WHERE 1 = 1
      AND od.order_id  = p_order_id
      AND i.user_id = @seller_id
      AND i.region_id = @region_id;

    -- seller_id
    -- dropshipper_id

    -- Update dropshipper inventory
    /*UPDATE
        inventory i
            INNER JOIN order_detail od
            ON i.`product_id` = od.`product_id`
            INNER JOIN `order` o
            ON od.order_id = o.id
    SET i.quantity = i.quantity - od.`quantity`
    WHERE 1 = 1
      AND od.order_id  = p_order_id
      AND i.user_id = @dropshipper_id
      AND i.region_id = @region_id;*/

    COMMIT;
END $$

DELIMITER;

