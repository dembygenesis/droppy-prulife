/**
 * Commands
 */

-- Update delivery_detail
UPDATE
  delivery_detail dd
  INNER JOIN product p
    ON 1 = 1
    AND dd.product_id = p.id 
    SET dd.`price_per_item_distributor` = p.`price_per_item_dropshipper`,
  		dd.`total_price_distributor` = dd.`quantity` * p.`price_per_item_dropshipper`;

-- Update delivery (amount_distributor totals)
UPDATE
  delivery d
  INNER JOIN (SELECT
  d.`id`,
  IF(SUM(dd.`total_price_distributor`) IS NULL, 0, SUM(dd.`total_price_distributor`)) AS total_price_distributor
FROM
  delivery d
  LEFT JOIN delivery_detail dd
    ON 1 = 1
    AND d.`id` = dd.`delivery_id`
WHERE 1 = 1
GROUP BY d.`id`) a 
    ON 1 = 1
      AND d.`id` = a.id
SET amount_distributor = a.total_price_distributor;

-- Get all users having unequal balances
SELECT
    u.id,
    u.`created_date`,
    CONCAT(u.`lastname`, ', ', u.`firstname`) AS username,
    ut.`coin_amount` AS coin_amount_aggregate,
    SUM(ct.`amount`) AS coin_amount_single_aggregate,
    IF (ABS(SUM(ct.`amount`)) = ABS(ut.`coin_amount`), 1, 0) AS equal,
    (COUNT(DISTINCT(ct.id))) AS coin_transactions
FROM
    `user` u
        INNER JOIN user_total ut
                   ON 1 = 1
                       AND u.id = ut.`user_id`
        LEFT JOIN coin_transaction ct
                  ON 1 = 1
                      AND u.`id` = ct.`user_id`
WHERE 1 = 1
  AND ct.`is_active` = 1
  AND u.`user_type_id` != (SELECT id FROM user_type WHERE `name` = 'Admin')
GROUP BY u.`id`
ORDER BY coin_transactions DESC;