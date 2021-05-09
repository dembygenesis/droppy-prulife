DROP PROCEDURE IF EXISTS  `update_delivery`;
DELIMITER $$

CREATE PROCEDURE `update_delivery` (
    p_user_id INTEGER,
    p_delivery_id INTEGER,
    p_delivery_status TEXT,
    p_tracking_number TEXT,
    p_void_or_reject_reason TEXT
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;  -- rollback any changes made in the transaction
            RESIGNAL;  -- raise again the sql exception to the caller
        END;

    START TRANSACTION;


    -- Validate delivery ID
    SELECT
        COUNT(d.id),
        d.id,
        d.created_by,
        do.name,
        ds.name,
        d.region_id,
        d.seller_id,
        d.dropshipper_id,
        d.amount_distributor,
        d.declared_amount,
        d.service_fee
    INTO
        @delivery_exists,
        @delivery_id,
        @user_id,
        @delivery_option,
        @delivery_status,
        @region_id,
        @seller_id,
        @dropshipper_id,
        @delivery_amount_distributor,
        @delivery_declared_amount,
        @delivery_service_fee
    FROM delivery d
             INNER JOIN delivery_option do
                        ON 1 = 1
                            AND do.id = d.delivery_option_id

             INNER JOIN delivery_status ds
                        ON 1 = 1
                            AND ds.id = d.delivery_status_id
    WHERE 1 = 1
      AND d.id = p_delivery_id
      AND d.is_active = 1;

    IF @delivery_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Delivery does not exist!';
    END IF;

    -- Validate delivery status
    SELECT
        COUNT(id),
        id
    INTO @delivery_status_exists, @delivery_status_id
    FROM delivery_status ds
    WHERE 1 = 1
        AND ds.name = p_delivery_status;

    IF @delivery_status_exists = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Delivery Status does not exist!';
    END IF;

    -- Validate user
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

    -- Seller logic
    IF @user_type = 'Seller' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Sellers cannot update deliveries. You smoking crack?';
    END IF;

    SET @admin_account = get_max_admin_id();
    SET @DROPSHIPPER_FEE = (SELECT CAST(value AS DECIMAL(65,2)) FROM sysparam WHERE `key` = 'DROPSHIPPER_FEE');

    -- Dropshipper logic
    IF @user_type = 'Dropshipper' THEN
        SELECT
            COUNT(d.id),
            ds.name
        INTO @has_delivery_access, @delivery_status
        FROM delivery d
            INNER JOIN delivery_status ds
        ON 1 = 1
            AND d.delivery_status_id = ds.id
        WHERE 1 = 1
          AND d.id = p_delivery_id
          AND d.dropshipper_id = p_user_id
          AND d.is_active = 1;

        IF @has_delivery_access = 0 THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Dropshippers can only update delivery entries they are assigned to.';
        END IF;

        /**
          "Accepted" and "Same Day Delivery" are handled the SAME way.
         */
        IF @delivery_status = 'Proposed' AND p_delivery_status IN ('Accepted', 'Same Day Delivery') THEN
            -- You can make it accepted.
            UPDATE delivery
            SET delivery_status_id = (SELECT id FROM delivery_status WHERE name = p_delivery_status)
            WHERE id = p_delivery_id;

            INSERT INTO delivery_tracking (delivery_id, delivery_status_id, last_updated, updated_by)
            VALUES (
                        p_delivery_id,
                        @delivery_status_id,
                        NOW(),
                        p_user_id
                   );
        ELSEIF @delivery_status = 'Same Day Delivery Claimed' AND p_delivery_status = 'Same Day Delivery Delivered' THEN
            -- Claimed can either be "Delivered" or "Issues"
            UPDATE delivery
            SET delivery_status_id = (SELECT id FROM delivery_status WHERE name = p_delivery_status)
            WHERE id = p_delivery_id;

            INSERT INTO delivery_tracking (delivery_id, delivery_status_id, last_updated, updated_by)
            VALUES (
                       p_delivery_id,
                       @delivery_status_id,
                       NOW(),
                       p_user_id
                   );
        ELSEIF @delivery_status = 'Same Day Delivery Claimed' AND p_delivery_status = 'Same Day Delivery Issues' THEN
            -- Claimed can either be "Delivered" or "Issues"
            UPDATE delivery
            SET delivery_status_id = (SELECT id FROM delivery_status WHERE name = p_delivery_status)
            WHERE id = p_delivery_id;

            INSERT INTO delivery_tracking (delivery_id, delivery_status_id, last_updated, updated_by)
            VALUES (
                       p_delivery_id,
                       @delivery_status_id,
                       NOW(),
                       p_user_id
                   );
        ELSEIF @delivery_status = 'Proposed' AND p_delivery_status = 'Rejected' THEN

            -- You can reject proposed items.
            UPDATE delivery
            SET delivery_status_id = (SELECT id FROM delivery_status WHERE name = 'Rejected'),
                last_updated = NOW()
            WHERE id = p_delivery_id;

            INSERT INTO delivery_tracking (delivery_id, delivery_status_id, last_updated, updated_by)
            VALUES (
                       p_delivery_id,
                       @delivery_status_id,
                       NOW(),
                       p_user_id
                   );

            -- The tedious part (updating the inventory and coins)

            -- Extract details
            -- , @region_id, @seller_id, @dropshipper_id

            -- Update inventory (increase items)
            UPDATE inventory i
            INNER JOIN delivery_detail dt
                ON 1 = 1
                    AND dt.product_id = i.product_id
            SET i.quantity = i.quantity + dt.quantity
            WHERE 1 = 1
                AND dt.delivery_id = p_delivery_id
                AND i.region_id = @region_id
                AND i.seller_id = @seller_id
                AND i.dropshipper_id = @dropshipper_id;

            -- Void inserted coins for droppy and admin
            UPDATE coin_transaction
            SET is_active = 0
            WHERE delivery_id = p_delivery_id;

            -- Update totals for admin and seller

            -- Extract admin and it his deductible
            SELECT
                ct.user_id,
                ct.amount
            INTO @admin_id, @admin_deductible
            FROM coin_transaction ct
                INNER JOIN `user` u
                           ON 1 = 1
                               AND ct.user_id = u.id
                INNER JOIN user_type ut
                           ON 1 = 1
                               AND u.user_type_id = ut.id
                               AND ut.name = 'Admin'
            WHERE 1 = 1
              AND ct.delivery_id = p_delivery_id;

            -- Extract seller and his deductible
            SELECT
                ct.user_id,
                ct.amount
            INTO @seller_id, @seller_collectible
            FROM coin_transaction ct
                     INNER JOIN `user` u
                                ON 1 = 1
                                    AND ct.user_id = u.id
                     INNER JOIN user_type ut
                                ON 1 = 1
                                    AND u.user_type_id = ut.id
                                    AND ut.name = 'Seller'
            WHERE 1 = 1
              AND ct.delivery_id = p_delivery_id;

            -- Extract dropshipper and it his deductible
            SELECT
                ct.user_id,
                ct.amount
            INTO @dropshipper_id, @dropshipper_deductible
            FROM coin_transaction ct
                INNER JOIN `user` u
                           ON 1 = 1
                               AND ct.user_id = u.id
                INNER JOIN user_type ut
                           ON 1 = 1
                               AND u.user_type_id = ut.id
                               AND ut.name = 'Dropshipper'
            WHERE 1 = 1
              AND ct.delivery_id = p_delivery_id;


            -- Update total coin balances for both dropshipper and admin and seller (deduct)
            UPDATE user_total
                SET coin_amount = coin_amount - ABS(@admin_deductible)
            WHERE user_id = @admin_id;

            UPDATE user_total
                SET coin_amount = coin_amount - ABS(@dropshipper_deductible)
            WHERE user_id = @dropshipper_id;

            -- Update seller coins total
            -- @seller_id, @seller_collectible
            UPDATE user_total
            SET coin_amount = coin_amount + ABS(@seller_collectible)
            WHERE user_id = @seller_id;

        ELSEIF @delivery_status = 'Accepted' AND p_delivery_status = 'Fulfilled' THEN
            -- You can make it fulfilled.
            IF p_tracking_number = '' THEN
                SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Tracking number cannot be empty if you are going to make an item fulfilled.';
            END IF;

            UPDATE delivery
            SET delivery_status_id = (SELECT id FROM delivery_status WHERE name = 'Fulfilled'),
                tracking_number = p_tracking_number
            WHERE id = p_delivery_id;

            INSERT INTO delivery_tracking (delivery_id, delivery_status_id, last_updated, updated_by)
            VALUES (
                       p_delivery_id,
                       @delivery_status_id,
                       NOW(),
                       p_user_id
                   );
        ELSE
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'This Dropshipper has no valid actions allowed.';
        END IF;
    END IF;

    IF @user_type = 'Admin' THEN

        -- @delivery_exists,
        /*@user_id, @delivery_option, @delivery_status, @region_id, @seller_id, @dropshipper_id,
        @delivery_amount_distributor, @delivery_declared_amount*/

        -- Validate admin access.
        /*SELECT
            ds.name,
            d.amount_distributor,
            d.declared_amount,
            d.dropshipper_id
        INTO @delivery_status, , @dropshipper_id
        FROM delivery d
        INNER JOIN delivery_status ds
            ON 1 = 1
                AND d.delivery_status_id = ds.id
        WHERE 1 = 1
          AND d.id = p_delivery_id
          AND d.dropshipper_id = p_user_id
          AND d.is_active = 1;*/

        /**
            @delivery_exists,
            @user_id,
            @delivery_option,
            @delivery_status,
            @region_id,
            @seller_id,
            @dropshipper_id,
            @delivery_amount_distributor,
            @delivery_declared_amount

         */

        IF @delivery_status = 'Fulfilled' AND p_delivery_status = 'Delivered' THEN

            SET @admin_id = get_max_admin_id();

            -- Set status to delivered
            UPDATE delivery
            SET delivery_status_id = @delivery_status_id
            WHERE id = p_delivery_id;

            -- Add into tracking
            INSERT INTO delivery_tracking (delivery_id, delivery_status_id, last_updated, updated_by)
            VALUES (
                       p_delivery_id,
                       @delivery_status_id,
                       NOW(),
                       @admin_id
                   );

            /**
              If dropship:
              1. Reimburse dropshipper with ** (Distributor's price) **
              2. Reimburse seller with ** (Declared amount - Distributor's price) **
             */

            IF @delivery_option = 'Dropship' THEN
                /**
                  Seller
                 */

                -- @delivery_amount_distributor, @delivery_declared_amount

                -- Add coin transaction
                INSERT INTO `coin_transaction` (
                    `created_by`,
                    `created_date`,
                    `is_active`,
                    `user_id`,
                    `type`,
                    `amount`,
                    `delivery_id`
                )
                VALUES
                (
                    @admin_id,
                    NOW(),
                    1,
                    @seller_id,
                    'C',
                    ABS(@delivery_declared_amount - @delivery_amount_distributor) * -1,
                    p_delivery_id
                );

                -- Update user totals
                UPDATE user_total
                    SET coin_amount = coin_amount + (@delivery_declared_amount - @delivery_amount_distributor)
                WHERE 1 = 1
                    AND user_id = @seller_id;

                /**
                  Dropshipper
                 */
                -- Add coin transaction
                INSERT INTO `coin_transaction` (
                    `created_by`,
                    `created_date`,
                    `is_active`,
                    `user_id`,
                    `type`,
                    `amount`,
                    `delivery_id`
                )
                VALUES
                (
                    @admin_id,
                    NOW(),
                    1,
                    @dropshipper_id,
                    'C',
                    ABS(@delivery_amount_distributor) * -1,
                    p_delivery_id
                );

                -- Update dropshipper totals
                UPDATE user_total
                    SET coin_amount = coin_amount + (@delivery_amount_distributor)
                WHERE 1 = 1
                  AND user_id = @dropshipper_id;

                /**
                  Admin
                 */
                -- Credit declared value

                -- Add coin transaction
                INSERT INTO `coin_transaction` (
                    `created_by`,
                    `created_date`,
                    `is_active`,
                    `user_id`,
                    `type`,
                    `amount`,
                    `delivery_id`
                )
                VALUES
                (
                    @admin_id,
                    NOW(),
                    1,
                    @admin_id,
                    'C',
                    ABS(@delivery_declared_amount) * -1,
                    p_delivery_id
                );

                -- Update admin totals
                UPDATE user_total
                SET coin_amount = coin_amount - (@delivery_declared_amount)
                WHERE 1 = 1
                  AND user_id = @admin_id;
            END IF;

            /**
              If Parcel:
              1. Reimburse seller with ** (Declared amount) **
             */

            IF @delivery_option = 'Parcel' THEN


                /**
                  Seller
                 */
                -- Add coin transaction
                INSERT INTO `coin_transaction` (
                    `created_by`,
                    `created_date`,
                    `is_active`,
                    `user_id`,
                    `type`,
                    `amount`,
                    `delivery_id`
                )
                VALUES
                (
                    @admin_id,
                    NOW(),
                    1,
                    @seller_id,
                    'C',
                    ABS(@delivery_declared_amount) * -1,
                    p_delivery_id
                );

                -- Update user totals
                UPDATE user_total
                    SET coin_amount = coin_amount + (@delivery_declared_amount)
                WHERE 1 = 1
                  AND user_id = @seller_id;

                /**
                  Admin
                 */
                -- Add coin transaction
                INSERT INTO `coin_transaction` (
                    `created_by`,
                    `created_date`,
                    `is_active`,
                    `user_id`,
                    `type`,
                    `amount`,
                    `delivery_id`
                )
                VALUES
                (
                    @admin_id,
                    NOW(),
                    1,
                    @admin_id,
                    'C',
                    ABS(@delivery_declared_amount) * -1,
                    p_delivery_id
                );

                -- Update admin totals
                UPDATE user_total
                SET coin_amount = coin_amount - (@delivery_declared_amount)
                WHERE 1 = 1
                  AND user_id = @admin_id;
            END IF;

        ELSEIF @delivery_status = 'Fulfilled' AND p_delivery_status = 'Returned' THEN

            IF p_void_or_reject_reason = '' THEN
                SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'No void or reject reason provided';
            END IF;

            UPDATE delivery
                SET delivery_status_id = (SELECT id FROM delivery_status WHERE name = 'Returned'),
                    void_or_reject_reason = p_void_or_reject_reason
            WHERE id = p_delivery_id;

            /**
              UPDATE THIS TOMORROW LOL
             */

            -- Void all coin transactions (just return the goddamn service fee and LEAVE THIS ALONE!!)
            -- Ok, this no longer applies
            /*UPDATE coin_transaction
                SET is_active = 0
            WHERE delivery_id = @delivery_id;*/

            -- Remove (service fee - @DROPSHIPPER_FEE) from admin
            /*UPDATE user_total
                SET coin_amount = coin_amount - (@delivery_service_fee - @DROPSHIPPER_FEE)
            WHERE user_id = @admin_account;*/

            -- Remove 35 from dropshipper
            /*UPDATE user_total
                SET coin_amount = coin_amount - @DROPSHIPPER_FEE
            WHERE user_id = @admin_account;*/

            -- Increase seller by service fee (we will reverse this later, for now just reverse everything)
            /*UPDATE user_total
                SET coin_amount = coin_amount + (@delivery_service_fee)
            WHERE user_id = @seller_id;*/

            IF @delivery_option = 'Parcel' THEN
                -- Return inventory
                UPDATE inventory i
                INNER JOIN delivery_detail dd
                    ON 1 = 1
                        AND dd.delivery_id = @delivery_id
                        AND i.region_id = @region_id
                        AND i.seller_id = @seller_id
                        AND i.dropshipper_id = @dropshipper_id
                        AND dd.product_id = i.product_id
                SET i.quantity = i.quantity + dd.quantity;
            END IF;

            /**
              The reject part lol
             */

            -- Add a coin balance to admin because TY nang service fee as a consequence of rejection
            /*INSERT INTO `coin_transaction` (
                `created_by`,
                `created_date`,
                `is_active`,
                `user_id`,
                `type`,
                `amount`,
                `delivery_id`
            )
            VALUES
            (
                @admin_id,
                NOW(),
                1,
                @admin_id,
                'D',
                @delivery_service_fee,
                p_delivery_id
            );*/

            -- Add a coin balance to seller because TY nang service fee as a consequence of rejection
            /*INSERT INTO `coin_transaction` (
                `created_by`,
                `created_date`,
                `is_active`,
                `user_id`,
                `type`,
                `amount`,
                `delivery_id`
            )
            VALUES
            (
                @admin_id,
                NOW(),
                1,
                @seller_id,
                'D',
                @delivery_service_fee,
                p_delivery_id
            );*/

            -- Update admin totals (increase service fee)
            /*UPDATE user_total
                SET coin_amount = coin_amount + @delivery_service_fee
            WHERE user_id = @admin_account;*/

            -- Update seller totals (decrease service fee)
            /*UPDATE user_total
                SET coin_amount = coin_amount - @delivery_service_fee
            WHERE user_id = @seller_id;*/

        ELSEIF @delivery_status != 'Voided' AND p_delivery_status = 'Voided' THEN

            -- Ensure validate void reason is present
            IF p_void_or_reject_reason = '' THEN
                SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'No void or reject reason provided';
            END IF;


            IF @delivery_status IN ('Returned') THEN
                /**
                  Logic for: Returned
                    * Void coin transactions
                    * Void delivery
                    * Remove (service fee) from admin
                    * Add (service_fee) added to seller
                 */

                -- Void coin transactions
                UPDATE coin_transaction
                    SET is_active = 0
                WHERE 1 = 1
                  AND delivery_id = p_delivery_id;

                -- Void delivery
                UPDATE delivery
                    SET is_active = 0,
                    delivery_status_id = (SELECT id FROM delivery_status WHERE `name` = 'Voided')
                WHERE 1 = 1
                  AND id = p_delivery_id;

                -- AAAAAAAAAAAAAAAAAAAA

                -- Decrease admin service fee share
                UPDATE user_total
                    SET coin_amount = coin_amount - (@delivery_service_fee - @DROPSHIPPER_FEE)
                WHERE user_id = @admin_account;

                -- Decrease the 35 service fee share from dropshipper
                UPDATE user_total
                    SET coin_amount = coin_amount - @DROPSHIPPER_FEE
                WHERE user_id = @admin_account;

                -- Increase seller by service fee charged
                UPDATE user_total
                    SET coin_amount = coin_amount + (@delivery_service_fee)
                WHERE user_id = @seller_id;

                -- AAAAAAAAAAAAAAAAAAAA


                -- Remove (service fee) from admin
                /*UPDATE user_total
                    SET coin_amount = coin_amount - @delivery_service_fee
                WHERE 1 = 1
                  AND user_id = @admin_account;*/

                -- Remove 35 (service fee portion) from dropshipper
                /*UPDATE user_total
                    SET coin_amount = coin_amount - @DROPSHIPPER_FEE
                WHERE 1 = 1
                  AND user_id = @dropshipper_id;*/

                -- Add (service_fee) to seller
                /*UPDATE user_total
                    SET coin_amount = coin_amount + @delivery_service_fee
                WHERE 1 = 1
                  AND user_id = @seller_id;*/
            END IF;

            /**
              Proposed, Accepted, and Fulfilled are handled THE SAME.

              While...

              Reject and Delivered are handled ANOTHER WAY.
             */

            IF @delivery_status IN (
                'Proposed',
                'Accepted',
                'Same Day Delivery',
                'Same Day Delivery Claimed',
                'Same Day Delivery Delivered',
                'Same Day Delivery Issues',
                'Fulfilled'
            ) THEN

                /**
                  Logic for: BOTH
                  * Void coin transactions
                  * Void delivery
                  * Remove (35 coins) from dropshipper
                  * Remove (service fee) from admin
                  * Add (service fee) to seller
                 */

                -- Void coin transactions
                UPDATE coin_transaction
                    SET is_active = 0
                WHERE 1 = 1
                  AND delivery_id = p_delivery_id;

                -- Void delivery
                UPDATE delivery
                    SET is_active = 0,
                        delivery_status_id = (SELECT id FROM delivery_status WHERE `name` = 'Voided')
                WHERE 1 = 1
                  AND id = p_delivery_id;

                -- Remove (service fee - @DROPSHIPPER_FEE) from admin
                UPDATE user_total
                    SET coin_amount = coin_amount - (@delivery_service_fee - @DROPSHIPPER_FEE)
                WHERE 1 = 1
                  AND user_id = @admin_account;

                -- Remove (35 coins) from dropshipper
                UPDATE user_total
                SET coin_amount = coin_amount - @DROPSHIPPER_FEE
                WHERE 1 = 1
                  AND user_id = @dropshipper_id;

                -- Add (service fee) back to seller
                UPDATE user_total
                    SET coin_amount = coin_amount + @delivery_service_fee
                WHERE 1 = 1
                  AND user_id = @seller_id;

                /**
                  !!! IF PARCEL UPDATE INVENTORY !!!
                 */
                IF @delivery_option = 'Parcel' THEN
                    UPDATE inventory i
                        INNER JOIN delivery_detail dt
                        ON 1 = 1
                            AND dt.product_id = i.product_id
                    SET i.quantity = i.quantity + dt.quantity
                    WHERE 1 = 1
                      AND dt.delivery_id = p_delivery_id
                      AND i.region_id = @region_id
                      AND i.seller_id = @seller_id
                      AND i.dropshipper_id = @dropshipper_id;
                END IF;
            END IF;

            /*IF @delivery_status IN ('Fulfilled') THEN
                SELECT 'I am king of wakanda';
            END IF;*/
        ELSE
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'This Admin has no valid actions allowed.';
        END IF;
    END IF;


    -- SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'It is the end.';

    COMMIT;
END $$

DELIMITER;

