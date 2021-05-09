/**
  Transactions
 */

DROP PROCEDURE IF EXISTS  `void_delivery`;
DELIMITER $$

DELIMITER;

CREATE PROCEDURE `void_delivery` (
    -- created by
    p_admin_id INTEGER,
    p_delivery_id INTEGER
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

    -- Validate delivery ID
    SELECT
        COUNT(d.id),
        d.created_by,
        do.name
    INTO @delivery_exists, @user_id, @delivery_option
    FROM delivery d
    INNER JOIN delivery_option do
        ON 1 = 1
            AND do.id = d.delivery_option_id
    WHERE 1 = 1
        AND d.id = p_delivery_id
        AND d.is_active = 1;

    IF @delivery_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Delivery does not exist!';
    END IF;

    -- Void delivery and reduce inventory
    UPDATE `delivery`
        SET is_active = 0
    WHERE id = p_delivery_id;

    -- Update inventory and reduce product quantity if this is a parcel entry
    IF @delivery_option = 'Parcel' THEN
        UPDATE inventory i
            INNER JOIN delivery_detail dd
            ON 1 = 1
                AND i.product_id = dd.product_id
        SET i.quantity = i.quantity + dd.quantity
        WHERE 1 = 1
          AND dd.delivery_id = p_delivery_id
          AND dd.product_id = i.product_id;
    END IF;

    COMMIT;
END $$

