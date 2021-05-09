/**
  Transactions
 */

DROP PROCEDURE IF EXISTS  `add_transaction`;
DELIMITER $$

CREATE PROCEDURE `add_transaction` (
    -- main params
    p_amount DECIMAL(15,2),
    p_coin_amount DECIMAL(15,2),
    p_admin_allotted_id INTEGER,
    p_user_allotted_id INTEGER,
    p_money_in BOOLEAN,
    p_bank_type_id INTEGER,
    p_reference_number TEXT,
    p_description TEXT
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
            AND user.id = p_admin_allotted_id
    );

    IF @admin_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Admin does not exist';
    END IF;


    -- Validate user (either can be dropshipper or seller)
    SET @seller_exists = (
        SELECT
            COUNT(user.id)
        FROM user
        INNER JOIN user_type
            ON 1 = 1
                AND user.user_type_id = user_type.id
        WHERE 1 = 1
          AND is_active = 1
          AND user_type.name IN ('Seller', 'Dropshipper')
          AND user.id = p_user_allotted_id
    );

    IF @seller_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Seller/Dropshipper does not exist';
    END IF;

    -- Validate amounts
    IF p_amount IS NULL OR p_amount < 0 THEN
        SET @message = CONCAT('Invalid amount provided, must not be NULL or negative: ', p_coin_amount);
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = @message;
    END IF;

    IF p_coin_amount IS NULL OR p_coin_amount < 0 THEN
        SET @message = CONCAT('Invalid coin_amount provided, must not be NULL or negative: ', p_coin_amount);
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = @message;
    END IF;

    -- Validate bank type id
    SET @bank_type_id_exists = (
        SELECT
            COUNT(*)
        FROM bank_type
        WHERE 1 = 1
            AND id = p_bank_type_id
    );

    IF @bank_type_id_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Bank Type ID does not exist';
    END IF;

    -- Convert amounts to either positive of negative depending on p_money_in
    IF p_money_in = true THEN
        SET p_amount = ABS(p_amount);
        SET p_coin_amount = ABS(p_coin_amount);
    ELSE
        SET p_amount = ABS(p_coin_amount) * - 1;
        SET p_coin_amount = ABS(p_coin_amount) * - 1;
    END IF;



    /**
      =========================
         **** MONEY IN ****
      =========================
     */

    -- Automatically populate the coins table if this is a money in
    IF p_money_in = true THEN
        -- Insert new transaction
        INSERT INTO `transaction` (
            `created_by`,
            `is_active`,
            `amount`,
            `coin_amount`,
            `admin_allotted_id`,
            `user_allotted_id`,
            `money_in`,
            `bank_type_id`,
            `reference_number`,
            `description`
        )
        VALUES
        (
            p_admin_allotted_id,
            1,
            p_amount,
            p_coin_amount,
            p_admin_allotted_id,
            p_user_allotted_id,
            p_money_in,
            p_bank_type_id,
            p_reference_number,
            p_description
        );

        -- Insert into coin transaction - log entry to admin
        SET @new_id = LAST_INSERT_ID();

        INSERT INTO coin_transaction (
            created_by,
            is_active,
            transaction_id,
            user_id,
            type,
            amount
        )
        VALUES
        (
            p_admin_allotted_id,
            1,
            @new_id,
            p_admin_allotted_id,
            'D', -- Debit
            p_coin_amount
        );

        -- Also, then insert another entry to a specific person
        INSERT INTO coin_transaction (
            created_by,
            is_active,
            user_id,
            type,
            amount,
            coin_transaction_id
        )
        VALUES
        (
            p_admin_allotted_id,
            1,
            p_user_allotted_id,
            'C', -- Debit
            ABS(p_coin_amount) * -1,
            LAST_INSERT_ID()
        );

        /**
          Update totals for both USER and ADMIN
         */

        -- Attempt to update user totals
        SET @user_total_exists = (
            SELECT
                COUNT(*)
            FROM `user_total`
            WHERE 1 = 1
                AND user_total.user_id = p_user_allotted_id
        );

        IF @user_total_exists = 1 THEN
            UPDATE `user_total`
                SET amount = amount + ABS(p_amount),
                    coin_amount = coin_amount + ABS(p_coin_amount)
            WHERE user_id = p_user_allotted_id;
        ELSE
            INSERT INTO `user_total` (user_id, amount, coin_amount, created_by, last_updated)
            VALUES (
                p_user_allotted_id,
                p_amount,
                p_coin_amount,
                p_admin_allotted_id,
                NOW()
            );
        END IF;

        -- Attempt to update admin totals
        SET @admin_total_exists = (
            SELECT
                COUNT(*)
            FROM `user_total`
            WHERE 1 = 1
              AND user_total.user_id = p_admin_allotted_id
        );

        IF @admin_total_exists = 1 THEN
            -- Update totals
            UPDATE `user_total`
            SET amount = amount + p_amount

                -- We don't update coins here because it totals out to zero anyway LOL,
                -- and we only update this when doing remittance? hmmm
                -- coin_amount = coin_amount + p_coin_amount
            WHERE user_id = p_admin_allotted_id;
        ELSE
            -- Make new entry
            INSERT INTO `user_total` (user_id, amount, coin_amount, created_by, last_updated)
            VALUES (
                       p_admin_allotted_id,
                       p_amount,
                       0,
                       p_admin_allotted_id,
                       NOW()
                   );
        END IF;

    ELSE
        /**
          =========================
             **** MONEY OUT ****
          =========================
         */

        -- The only scenario for money out right now is doing payments to Maximum 88 OR
        -- to sellers.

        -- So there are 2 instances of money out and different logic to handle them...
        -- Perhaps I should add another parameter?

        SELECT "Do money out";
    END IF;

    COMMIT;
END $$

DELIMITER;

