/**
  Transactions
 */

DROP PROCEDURE IF EXISTS  `void_transaction`;
DELIMITER $$

CREATE PROCEDURE `void_transaction` (
    -- main params
    p_transaction_id INTEGER,
    p_admin_id INTEGER
)
BEGIN

    -- Magic ;)
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
    );

    IF @admin_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Only an admin can delete transactions';
    END IF;

    -- Declare these variables to be used in validation below
    SET @transaction_exists = '';
    SET @amount = '';
    SET @coin_amount = '';
    SET @admin_allotted_id = '';
    SET @user_allotted_id = '';

    -- Validate transaction ID
    SELECT
        COUNT(id), amount, coin_amount, admin_allotted_id, user_allotted_id
    INTO @transaction_exists, @amount, @coin_amount, @admin_allotted_id, @user_allotted_id
    FROM `transaction`
    WHERE 1 = 1
      AND id = p_transaction_id
      AND is_active = 1;

    IF @transaction_exists = 0 THEN
        SET @message = CONCAT('Transaction ID: ', p_transaction_id, ' does not exist.');
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = @message;
    END IF;

    -- Attempt to void the transaction.
    UPDATE `transaction`
        SET is_active = 0
    WHERE id = p_transaction_id;

    -- Attempt to void the user totals for admin
    UPDATE user_total
        SET amount = amount - ABS(@amount)
    WHERE user_id = @admin_allotted_id;

    -- Attempt to void the user totals for seller
    UPDATE user_total
        SET amount = amount - ABS(@amount)
    WHERE user_id = @user_allotted_id;

    -- Also void the coin entry that came along with it
    UPDATE coin_transaction
        SET is_active = 0
    WHERE transaction_id = p_transaction_id;

    SET @coin_transaction_assigned_to_user_id = (
        SELECT
            id
        FROM coin_transaction
        WHERE 1 = 1
          AND transaction_id = p_transaction_id
    );

    -- Also void the coin entry assigned to the seller
    UPDATE coin_transaction
        SET is_active = 0
    WHERE coin_transaction_id = @coin_transaction_assigned_to_user_id;

    -- Update totals
    -- These are variables available to you: @transaction_exists, @amount, @coin_amount, @admin_allotted_id, @user_allotted_id
    -- Circumstantial voiding

    /**
      General Algo
      1. Void transaction ID
      2. Also void coin assignments related to that, and reflect to totals
      3. Also void coin assignments forwarded to users if it exists, and reflect to totals

      So strategy is, IF EXISTS - void and upate user totals
     */

    -- Void transaction ID
    UPDATE transaction
        SET is_active = 0,
            updated_by = p_admin_id
    WHERE id = p_transaction_id;

    -- Also void coin assignments related to that, and reflect to totals
    SET @coin_amount = '';
    SET @admin_id = '';
    SET @id = '';
    SET @count = '';
    SET @user_id = '';

    SELECT
        COUNT(id), amount, user_id, id
        INTO @count, @coin_amount, @admin_id, @id
    FROM coin_transaction
    WHERE 1 = 1
        AND transaction_id = p_transaction_id;

    -- If coin assignments exist
    IF @count != 0 THEN
        -- Update coin_transaction table
        UPDATE coin_transaction
            SET is_active = 0,
                updated_by = p_admin_id
        WHERE id = @id;

        UPDATE coin_transaction
        SET is_active = 0,
            updated_by = p_admin_id
        WHERE coin_transaction_id = @id;


        -- Update totals for admin ? -- NAH

        -- Update user totals for user? - Definitely
        SELECT
            user_id
        INTO @user_id
        FROM coin_transaction
        WHERE 1 = 1
            AND coin_transaction_id = @id;

        UPDATE user_total
            SET coin_amount = coin_amount - @coin_amount,
                updated_by = p_admin_id
        WHERE user_id = @user_id;
    END IF;


    COMMIT;
END $$

DELIMITER;

