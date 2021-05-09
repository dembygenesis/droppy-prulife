

DELIMITER $$
DROP PROCEDURE IF EXISTS `add_test`$$
CREATE PROCEDURE `add_test`()
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;
            -- rollback any changes made in the transaction
            RESIGNAL;
            -- raise again the sql exception to the caller
        END;
    START TRANSACTION;

    INSERT INTO `test` (`ss`)
    VALUES
    ('ss');

    CALL add_test_exception();
    COMMIT;
END$$
DELIMITER ;