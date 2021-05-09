DROP FUNCTION IF EXISTS  `get_base_price`;
DELIMITER $$

CREATE FUNCTION `get_base_price` (
    p_declared_amount INTEGER
) RETURNS INTEGER DETERMINISTIC
BEGIN
    DECLARE base_price INTEGER;

    IF p_declared_amount BETWEEN 0 AND 1499 THEN
        SET base_price = 130;
    ELSEIF p_declared_amount BETWEEN 1500 AND 1999 THEN
        SET base_price = 135;
    ELSEIF p_declared_amount BETWEEN 2000 AND 2499 THEN
        SET base_price = 135;
    ELSEIF p_declared_amount BETWEEN 2500 AND 2999 THEN
        SET base_price = 175;
    ELSEIF p_declared_amount BETWEEN 3000 AND 3499 THEN
        SET base_price = 180;
    ELSEIF p_declared_amount BETWEEN 3500 AND 3999 THEN
        SET base_price = 185;
    ELSEIF p_declared_amount BETWEEN 4000 AND 4499 THEN
        SET base_price = 190;
    ELSEIF (p_declared_amount BETWEEN 4500 AND 4999) OR p_declared_amount > 4999 THEN
        SET base_price = 195;
    END IF;

    RETURN base_price;
END $$

DELIMITER;

