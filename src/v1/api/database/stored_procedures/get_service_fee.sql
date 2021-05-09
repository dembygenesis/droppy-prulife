DROP PROCEDURE IF EXISTS  `get_service_fee`;
DELIMITER $$

CREATE PROCEDURE `get_service_fee` (
    p_declared_amount TEXT
)
this_proc:BEGIN
    DECLARE service_fee DECIMAL;

    IF p_declared_amount BETWEEN 0 AND 1499 THEN
        SET service_fee = 195;
    ELSEIF p_declared_amount BETWEEN 1500 AND 1999 THEN
        SET service_fee = 200;
    ELSEIF p_declared_amount BETWEEN 2000 AND 2499 THEN
        SET service_fee = 205;
    ELSEIF p_declared_amount BETWEEN 2500 AND 2999 THEN
        SET service_fee = 240;
    ELSEIF p_declared_amount BETWEEN 3000 AND 3499 THEN
        SET service_fee = 245;
    ELSEIF p_declared_amount BETWEEN 3500 AND 3999 THEN
        SET service_fee = 250;
    ELSEIF p_declared_amount BETWEEN 4000 AND 4499 THEN
        SET service_fee = 255;
    ELSEIF (p_declared_amount BETWEEN 4500 AND 4999) OR p_declared_amount > 4999 THEN
        SET service_fee = 260;
    END IF;

    SELECT service_fee;

    COMMIT;
END $$

DELIMITER;

