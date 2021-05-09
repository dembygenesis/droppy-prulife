DELIMITER $$
DROP PROCEDURE IF EXISTS `validate_user`$$
CREATE PROCEDURE `validate_user`(
    p_user_id INTEGER,
    p_role TEXT
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;
            -- rollback any changes made in the transaction
            RESIGNAL;
            -- raise again the sql exception to the caller
        END;

    -- Validate user query
    SELECT
        COUNT(u.id),
        ut.name
    INTO @user_exists, @user_type
    FROM `user` u
    INNER JOIN user_type ut
        ON 1 = 1
            AND ut.id = u.user_type_id
    WHERE 1 = 1
      AND u.is_active = 1
      AND u.id = p_user_id;

    IF @user_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'User does not exist';
    END IF;

    IF @user_type != p_role THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'User role mismatch';
    END IF;
END$$
DELIMITER ;