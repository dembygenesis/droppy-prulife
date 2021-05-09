/**
  Reject
  Void
  Accept
 */


DROP PROCEDURE IF EXISTS  `update_withdrawal`;
DELIMITER $$

CREATE PROCEDURE `update_withdrawal` (
    p_user_id INTEGER,
    p_withdrawal_id INTEGER,
    p_withdrawal_status TEXT,
    p_reference_number TEXT,
    p_void_or_reject_reason TEXT
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;  -- rollback any changes made in the transaction
            RESIGNAL;  -- raise again the sql exception to the caller
        END;

    START TRANSACTION;

    -- Validate admin (as they are only the ones who can perform actions)
    SELECT
        COUNT(id),
        id
    INTO @user_exists, @user_id
    FROM `user` u
    WHERE 1 = 1
        AND u.id = p_user_id
        AND u.is_active = 1
        AND u.user_type_id = (SELECT id FROM user_type WHERE `name` = 'Admin')
    ;

    IF @user_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Admin does not exist';
    END IF;

    -- Validate status
    SELECT
        COUNT(id),
        id
    INTO @withdrawal_status_exists, @withdrawal_status_id
    FROM withdrawal_status wt
    WHERE 1 = 1
        AND wt.name = p_withdrawal_status;

    IF @withdrawal_status_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Withdrawal status does not exist';
    END IF;

    -- Validate withdrawal
    SELECT
        COUNT(w.id),
        ws.name,
        w.total_amount,
        w.fee,
        w.user_id,
        w.bank_type_id
    INTO
        @withdrawal_status_exists,
        @withdrawal_current_status,
        @withdrawal_total_amount,
        @withdrawal_fee,
        @withdrawal_user_id,
        @withdrawal_bank_type_id
    FROM withdrawal w
    INNER JOIN withdrawal_status ws
        ON 1 = 1
            AND w.withdrawal_status_id = ws.id
    WHERE 1 = 1
      AND w.id = p_withdrawal_id
      AND w.is_active = 1
    ;

    SET @admin_account = get_max_admin_id();

    IF @withdrawal_status_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Withdrawal does not exist';
    END IF;

    /**
      Available actions
     */

    IF @withdrawal_current_status = 'Pending' AND p_withdrawal_status = 'Completed' THEN

        -- Validate reference number
        IF p_reference_number = '' THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Reference Number needed';
        END IF;

        -- Accept withdrawal
        UPDATE withdrawal
        SET withdrawal_status_id = @withdrawal_status_id,
            reference_number = p_reference_number,
            last_updated = NOW(),
            updated_by = @admin_account
        WHERE 1 = 1
          AND id = p_withdrawal_id;

        -- Add CREDIT coin transaction to "Admin"
        INSERT INTO coin_transaction (created_by, created_date, is_active, user_id, type, amount, withdrawal_id)
        VALUES (
                   @admin_account,
                   NOW(),
                   1,
                   @admin_account,
                   'C',
                   ABS(@withdrawal_total_amount) * -1,
                   p_withdrawal_id
               );

        -- Add CREDIT transaction to "Admin"
        INSERT INTO `transaction` (
            created_by,
            created_date,
            is_active,
            amount,
            coin_amount,
            admin_allotted_id,
            user_allotted_id,
            money_in,
            bank_type_id,
            reference_number,
            description,
            withdrawal_id
        )
        VALUES
        (
            @admin_account,
            NOW(),
            1,
            (@withdrawal_total_amount - @withdrawal_fee) * -1,
            0,
            @admin_account,
            @admin_account,
            0,
            @withdrawal_bank_type_id,
            p_reference_number,
            CONCAT('For: ', p_reference_number),
            p_withdrawal_id
        );

        -- Add CREDIT transaction to "User"
        INSERT INTO `transaction` (
            created_by,
            created_date,
            is_active,
            amount,
            coin_amount,
            admin_allotted_id,
            user_allotted_id,
            money_in,
            bank_type_id,
            reference_number,
            description,
            withdrawal_id
        )
        VALUES
        (
            @admin_account,
            NOW(),
            1,
            (@withdrawal_total_amount - @withdrawal_fee),
            0,
            @admin_account,
            @withdrawal_user_id,
            1,
            @withdrawal_bank_type_id,
            p_reference_number,
            CONCAT('For withdrawal: ', p_withdrawal_id),
            p_withdrawal_id
        );

        -- Update coin totals for admin
        UPDATE user_total
            SET coin_amount = coin_amount - @withdrawal_total_amount
        WHERE 1 = 1
            AND user_id = @admin_account
        ;

        -- Update money totals for admin
        UPDATE user_total
            SET amount = amount - (@withdrawal_total_amount - @withdrawal_fee)
        WHERE 1 = 1
          AND user_id = @admin_account
        ;

        -- Update money Totals for seller
        UPDATE user_total
            SET amount = amount + (@withdrawal_total_amount - @withdrawal_fee)
        WHERE 1 = 1
          AND user_id = @withdrawal_user_id
        ;

    ELSEIF (@withdrawal_current_status = 'Completed' || @withdrawal_current_status = 'Pending') AND p_withdrawal_status = 'Voided' THEN
        -- Validate void rejection
        IF p_void_or_reject_reason = '' THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Reference Number needed';
        END IF;

        -- Update withdrawal
        UPDATE withdrawal
        SET withdrawal_status_id = @withdrawal_status_id,
            void_or_reject_reason = p_void_or_reject_reason,
            last_updated = NOW(),
            updated_by = @admin_account
        WHERE 1 = 1
          AND id = p_withdrawal_id;

        /**
          Logic for rejecting "Pending" items.
         */
        IF @withdrawal_current_status = 'Pending' THEN
            -- Void the coin transactions
            UPDATE coin_transaction
                SET is_active = 0
            WHERE withdrawal_id = p_withdrawal_id;

            -- Reverse the coin totals for user (seller/dropshiper) (INCREASE: Return coins)
            UPDATE user_total
                SET coin_amount = coin_amount + @withdrawal_total_amount
            WHERE user_id = @withdrawal_user_id;

            -- Reverse the coin totals for admin (DECREASE)
            UPDATE user_total
                SET coin_amount = coin_amount - @withdrawal_total_amount
            WHERE user_id = @admin_account;
        END IF;

        /**
          Logic for rejecting "Completed" items.
         */
        IF @withdrawal_current_status = 'Completed' THEN
            -- Void the coin transactions
            UPDATE coin_transaction
                SET is_active = 0
            WHERE withdrawal_id = p_withdrawal_id;

            -- Void the transactions
            UPDATE `transaction`
                SET is_active = 0
            WHERE withdrawal_id = p_withdrawal_id;

            -- Add user coins back
            UPDATE user_total
                SET coin_amount = coin_amount + @withdrawal_total_amount
            WHERE user_id = @withdrawal_user_id;

            -- Add admin coins back
            UPDATE user_total
            SET coin_amount = coin_amount + @withdrawal_total_amount
            WHERE user_id = @admin_account;

            -- Add admin money
            UPDATE user_total
                SET amount = amount + (@withdrawal_total_amount - @withdrawal_fee)
            WHERE user_id = @admin_account;

            -- Decrease user money
            UPDATE user_total
                SET amount = amount - (@withdrawal_total_amount - @withdrawal_fee)
            WHERE user_id = @withdrawal_user_id;
        END IF;

        /**
          !!! COMPLETED !!!
          Apply different logic for voiding completed items (void withdrawals, coin transactions, money transactions
          and return totals
         */

        /**
          !!! Pending !!!
          Apply different logic for voiding completed items (void withdrawals, coin transactions, money transactions
          and return totals
         */

    ELSE
        -- Raise error
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'No available actions';
    END IF;

    COMMIT;
END $$

DELIMITER;

