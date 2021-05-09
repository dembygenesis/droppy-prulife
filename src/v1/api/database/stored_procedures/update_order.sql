DROP PROCEDURE IF EXISTS  `update_order`;
DELIMITER $$

/**
  This "update_order" command only serves tas deleting order transactions
 */

CREATE PROCEDURE `update_order` (
    p_user_id INTEGER,
    p_order_id INTEGER,
    p_void_or_reject_reason TEXT
)
BEGIN

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;  -- rollback any changes made in the transaction
            RESIGNAL;  -- raise again the sql exception to the caller
        END;

    START TRANSACTION;

    IF p_void_or_reject_reason = '' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Void reason empty';
    END IF;

    -- Validate default admin.
    SET @admin_exists = (
        SELECT COUNT(id) FROM `user`
        WHERE 1 = 1
          AND is_active = 1
          AND id = p_user_id
          AND user_type_id = (SELECT `id` FROM user_type WHERE `name` = 'Admin')
    );

    IF @admin_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Admin does not exist';
    END IF;

    SET @order_status_id_voided = (
        SELECT
            id
        FROM order_status
        WHERE 1 = 1
          AND `name` = 'Voided'
    );

    SET @order_status_id_proposed = (
        SELECT
            id
        FROM order_status
        WHERE 1 = 1
          AND `name` = 'Proposed'
    );

    -- Validate order
    SELECT
        COUNT(id),
        id,
        amount,
        region_id,
        seller_id,
        dropshipper_id
    INTO @order_count, @order_id, @order_amount, @order_region_id, @order_seller_id, @order_dropshipper_id
    FROM `order` o
    WHERE 1 = 1
        AND o.is_active = 1
        AND o.order_status_id = @order_status_id_proposed
        AND o.id = p_order_id;

    IF @order_count != 1 THEN
        -- SET @message = CONCAT('order count, ', @order_count, ' hahaha @order_status_id_proposed', @order_status_id_proposed, ' p_order_id ', p_order_id);
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Order does not exist.';
        -- SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = @message;
    END IF;

    -- Void order
    UPDATE `order`
        SET is_active = 0,
            order_status_id = @order_status_id_voided,
            void_or_reject_reason = p_void_or_reject_reason,
            updated_by = p_user_id,
            last_updated = NOW()
    WHERE id = @order_id;

    -- Void order coin_transactions
    UPDATE coin_transaction
        SET is_active = 0,
            updated_by = p_user_id,
            last_updated = NOW()
    WHERE order_id = @order_id;

    /*SELECT *, 'pre_inventory' AS typee FROM inventory i
    WHERE 1 = 1
      AND i.dropshipper_id = @order_dropshipper_id
      AND i.seller_id = @order_seller_id
      AND i.region_id = @order_region_id
    ;

    SELECT * FROM order_detail od
    WHERE 1 = 1
      AND od.order_id = @order_id
    ;*/

    -- Update user inventory
    UPDATE inventory i
    INNER JOIN order_detail od
        ON 1 = 1
            AND od.order_id = @order_id
            AND i.region_id = @order_region_id
            AND i.product_id = od.product_id
            AND i.dropshipper_id = @order_dropshipper_id
            AND i.seller_id = @order_seller_id
    SET
        i.quantity = i.quantity - od.quantity;

    -- Check if there are negative balances, and FAIL if there are
    SET @negative_inventory_count = (
        SELECT
            COUNT(*)
        FROM inventory i
        WHERE 1 = 1
          AND i.region_id = @order_region_id
          AND i.seller_id = @order_seller_id
          AND i.dropshipper_id = @order_dropshipper_id
          AND i.quantity < 0
    );

    IF @negative_inventory_count != 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'You cannot void orders if they result in negative product inventories.';
    END IF;

    SELECT @negative_inventory_count;

    /*SELECT @order_count, @order_id, @order_amount, @order_region_id, @order_seller_id, @order_dropshipper_id;

    SELECT *, 'post_inventory' AS typee FROM inventory i
    WHERE 1 = 1
      AND i.dropshipper_id = @order_dropshipper_id
      AND i.seller_id = @order_seller_id
      AND i.region_id = @order_region_id
    ;*/

    -- Update coins for dropshipper -> reduce by 75
    UPDATE user_total
        SET coin_amount = coin_amount - 75
    WHERE user_id = @order_dropshipper_id;

    -- Update coins for seller -> add the amount (10,750) as of now
    UPDATE user_total
        SET coin_amount = coin_amount + @order_amount
    WHERE user_id = @order_seller_id;

    -- Update coins for admin -> reduce the amount by (10,750) - 75
    SET @admin_account = get_max_admin_id();

    UPDATE user_total
        SET coin_amount = coin_amount - (@order_amount - 75)
    WHERE user_id = @admin_account;

    -- SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'End TEST';
    COMMIT;
END $$

DELIMITER;

