DROP FUNCTION IF EXISTS  `get_max_admin_id`;
DELIMITER $$

CREATE FUNCTION `get_max_admin_id` (

) RETURNS INTEGER DETERMINISTIC
BEGIN
    SELECT
        MAX(u.id)
    INTO @admin_id
    FROM `user` u
    INNER JOIN user_type ut
        ON 1 = 1
            AND u.user_type_id = ut.id
    WHERE 1 = 1
        AND ut.name = 'Admin'
    ;

    RETURN @admin_id;
END $$

DELIMITER;

