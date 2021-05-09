/**
  Transactions
 */

DROP PROCEDURE IF EXISTS  `add_coin_transaction`;
DELIMITER $$

CREATE PROCEDURE `add_coin_transaction` (
    -- main params
    p_amount DECIMAL(15,2),
    p_coin_amount DECIMAL(15,2),
    p_admin_allotted_id INTEGER,
    p_user_allotted_id INTEGER,
    p_money_in BOOLEAN
)
BEGIN
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
          AND user.id = p_admin_allotted_id
    );

    IF @admin_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Admin does not exist';
    END IF;




    COMMIT;
END $$

DELIMITER;

