DELIMITER $$

USE `your_droppy_database`$$

DROP FUNCTION IF EXISTS `get_cumulative_balance`$$

CREATE FUNCTION `get_cumulative_balance`(
    p_user_id INTEGER
) RETURNS INT
    DETERMINISTIC
BEGIN
    SET @running_total := 0;

    SET @cumulative_balance = (
        SELECT
            cumulative_sum
        FROM
            (SELECT
                 ct.`id`,
                 ct.`amount`,
                 (
                     @running_total := @running_total + ct.amount
                     ) AS cumulative_sum
             FROM
                 coin_transaction ct
             WHERE 1 = 1
               AND ct.`user_id` = p_user_id
               AND ct.`is_active` = 1
             ORDER BY ct.`id` ASC) AS a
        ORDER BY id DESC LIMIT 1
    );

    RETURN @cumulative_balance;
END$$

DELIMITER ;