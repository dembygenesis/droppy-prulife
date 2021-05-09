/**
  Transactions
 */

DROP PROCEDURE IF EXISTS  `add_order`;
DELIMITER $$

CREATE PROCEDURE `add_order` (
    -- person authentication first
    p_admin_id INTEGER,
    p_user_id INTEGER,
    p_order_details TEXT
)
BEGIN

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;  -- rollback any changes made in the transaction
        RESIGNAL;  -- raise again the sql exception to the caller
    END;

    /*DECLARE `_rollback` BOOL DEFAULT 0;
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION SET `_rollback` = 1;

    SET AUTOCOMMIT = 0;*/
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

    INSERT INTO aaa
    SELECT 'abd';

    IF @admin_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Admin does not exist';
    END IF;

    COMMIT;
END $$

DELIMITER;

