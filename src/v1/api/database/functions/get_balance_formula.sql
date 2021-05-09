DELIMITER $$

$$

DROP FUNCTION IF EXISTS `get_balance_formula`$$

CREATE FUNCTION `get_balance_formula`(
    p_user_id INTEGER
) RETURNS INT
    DETERMINISTIC
BEGIN
    SET @running_total = 0;
    SET @user_id = p_user_id;

    SET @result = (SELECT
                       running_balance
                   FROM
                       (SELECT
                            *
                        FROM
                            (SELECT
                                 id,
                                 date_created,
                                 `type`,
                                 @running_total := @running_total + amount AS running_balance,
                                 amount,
                                 orig_amount,
                                 reference_number,
                                 bank_type,
                                 source,
                                 recipient,
                                 tran_num,
                                 is_active
                             FROM
                                 (
                                     (SELECT
                                          ct.id,
                                          IF(
                                                  DATE_FORMAT(
                                                          CONVERT_TZ(
                                                                  ct.created_date,
                                                                  '+00:00',
                                                                  '+08:00'
                                                              ),
                                                          '%Y-%m-%d %h:%i %p'
                                                      ) IS NULL,
                                                  "",
                                                  DATE_FORMAT(
                                                          CONVERT_TZ(
                                                                  ct.created_date,
                                                                  '+00:00',
                                                                  '+08:00'
                                                              ),
                                                          '%Y-%m-%d %h:%i %p'
                                                      )
                                              ) AS date_created,
                                          IF(
                                                      ct.type = 'C',
                                                      'Coins In',
                                                      'Coins Out'
                                              ) AS TYPE,
                                          ct.created_date,
                                          -- Add column here for types
                                          CASE
                                              WHEN ct.delivery_id IS NOT NULL
                                                  AND ds.name IN ('Rejected')
                                                  THEN 0
                                              ELSE ct.amount * - 1
                                              END AS amount,
                                          ct.amount AS orig_amount,
                                          -- End of type amounts
                                          -- ct.amount * - 1 AS amount,
                                          'N/A' AS reference_number,
                                          'N/A' AS bank_type,
                                          CASE
                                              WHEN ct.withdrawal_id IS NOT NULL
                                                  THEN 'Withdrawal'
                                              WHEN ct.delivery_id IS NOT NULL
                                                  THEN 'Delivery'
                                              WHEN ct.order_id IS NOT NULL
                                                  THEN 'Order'
                                              WHEN ct.coin_transaction_id IS NOT NULL
                                                  THEN 'Coins added from Cash In'
                                              END AS source,
                                          IF (d.id IS NULL, 'N/A', d.name) AS recipient,
                                          CASE
                                              WHEN ct.withdrawal_id IS NOT NULL
                                                  THEN ct.withdrawal_id
                                              WHEN ct.delivery_id IS NOT NULL
                                                  THEN ct.delivery_id
                                              WHEN ct.order_id IS NOT NULL
                                                  THEN ct.order_id
                                              WHEN ct.coin_transaction_id IS NOT NULL
                                                  THEN ct.coin_transaction_id
                                              END AS tran_num,
                                          ct.is_active
                                      FROM
                                          coin_transaction ct
                                              LEFT JOIN delivery d
                                                        ON 1 = 1
                                                            AND ct.delivery_id = d.id
                                              LEFT JOIN delivery_status ds
                                                        ON 1 = 1
                                                            AND d.delivery_status_id = ds.id
                                      WHERE 1 = 1
                                        AND ct.user_id = @user_id
                                        AND ct.is_active = 1
                                      ORDER BY ct.created_date DESC)
                                 ) AS a
                             ORDER BY created_date ASC) AS a
                        ORDER BY id DESC) AS a
                   LIMIT 1);

    RETURN @result;
END$$

DELIMITER ;