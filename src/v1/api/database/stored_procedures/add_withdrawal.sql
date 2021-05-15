/**
  Withdrawal
 */

DROP PROCEDURE IF EXISTS `add_withdrawal`;
DELIMITER $$

/**
  Dropshippers and users can request this
 */
CREATE PROCEDURE `add_withdrawal` (
    p_user_id INTEGER,
    p_amount INTEGER,
    p_bank_no TEXT,
    p_bank_type_id TEXT,
    p_bank_acount_name TEXT,
    p_contact_no TEXT
)
BEGIN

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;  -- rollback any changes made in the transaction
            RESIGNAL;  -- raise again the sql exception to the caller
        END;

    START TRANSACTION;

    -- Validate users
    SELECT
        COUNT(u.id),
        u.id,
        u_t.coin_amount
    INTO @user_count, @user_id, @coin_amount
    FROM `user` u
             INNER JOIN user_type ut
                        ON 1 = 1
                            AND u.user_type_id = ut.id
             INNER JOIN user_total u_t
                        ON 1 = 1
                            AND u.id = u_t.user_id
    WHERE 1 = 1
      AND u.id = p_user_id
      AND u.is_active = 1
      AND ut.name IN ('Seller', 'Dropshipper');

    -- Validate user
    IF @user_count = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Seller/Dropshipper does not exist';
    END IF;

    -- I treat constants as capital letters
    SET @WITHDRAWAL_FEE = (SELECT CAST(value AS DECIMAL(65,2)) FROM sysparam WHERE `key` = 'WITHDRAWAL_FEE');

    SET @amount = p_amount;
    SET @total_amount = p_amount;

    IF @amount = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Coins must not be 0.';
    END IF;

    IF @coin_amount < @total_amount THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'User does not have enough coins to make a withdrawal';
    END IF;

    IF @total_amount < @WITHDRAWAL_FEE THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'You must withdraw an amount greater than the withdrawal fee';
    END IF;

    SET @admin_account = get_max_admin_id();

    /**
      Add withdrawal
     */
    INSERT INTO `withdrawal` (
        `created_by`,
        `created_date`,
        `is_active`,
        `amount`,
        `fee`,
        `total_amount`,
        `user_id`,
        `withdrawal_status_id`,
        `bank_no`,
        `bank_type_id`,
        `bank_account_name`,
        `contact_no`
    )
    VALUES
    (
        @admin_account,
        NOW(),
        1,
        @amount,
        @WITHDRAWAL_FEE,
        @total_amount,
        p_user_id,
        (SELECT id FROM withdrawal_status WHERE `name` = 'Pending'),
        p_bank_no,
        p_bank_type_id,
        p_bank_acount_name,
        p_contact_no
    );

    SET @withdrawal_id_created = LAST_INSERT_ID();

    /**
      Subtract coins, and update totals
     */

    -- For user (debit for us) because they're giving us money
    INSERT INTO coin_transaction (created_by, created_date, is_active, user_id, type, amount, withdrawal_id)
    VALUES (
               @admin_account,
               NOW(),
               1,
               p_user_id,
               'D',
               @total_amount,
               @withdrawal_id_created
           );

    -- For admin
    INSERT INTO coin_transaction (created_by, created_date, is_active, user_id, type, amount, withdrawal_id)
    VALUES (
               @admin_account,
               NOW(),
               1,
               @admin_account,
               'D',
               @total_amount,
               @withdrawal_id_created
           );

    /**
      Update totals
     */

    -- Update user
    UPDATE user_total
    SET coin_amount = coin_amount - ABS(@total_amount)
    WHERE user_id = p_user_id;

    -- Update admin
    UPDATE user_total
    SET coin_amount = coin_amount + ABS(@total_amount)
    WHERE user_id = @admin_account;

    -- SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Set delivery test';

    COMMIT;
END $$

DELIMITER;

