/**
 	IMPORTANT: 
 	This migration script currently does NOT consider voided items,
 	please keep that in mind when running it.

 	Eventually as the app grows and there are some mistakes, this week need to be re-ran again.
 */



/**
 * Delete all coins not stemming from manual processes (coin transactions or transactions)
 */

DELETE
  e.*
FROM
  coin_transaction e
WHERE id NOT IN
  (SELECT
    id
  FROM
    (SELECT
      id
    FROM
      coin_transaction ct
    WHERE ct.`transaction_id` IS NOT NULL
      OR ct.`coin_transaction_id` IS NOT NULL) X);

/**
 * Delete User Totals
 */
DELETE FROM `user_total`;


/**
 * Repopulate User Total
 */

INSERT INTO user_total (
  user_id,
  amount,
  coin_amount,
  created_by
)
(SELECT
  id,
  0,
  0,
  (SELECT
    MAX(u.id)
  FROM
    `user` u
    INNER JOIN user_type ut
      ON 1 = 1
      AND u.user_type_id = ut.id
      AND ut.name = 'Admin')
FROM
  `user`);

/**
 * Repopulate user total cash for sellers
 */

UPDATE
  user_total ut
  INNER JOIN
    (SELECT
      t.`user_allotted_id`,
      SUM(t.`amount`) AS amount
    FROM
      `transaction` t
    GROUP BY t.`user_allotted_id`) a
    ON 1 = 1
    AND ut.`user_id` = a.user_allotted_id 
    SET ut.`amount` = a.amount
WHERE ut.`user_id` = a.user_allotted_id;

/**
 * Repopulate user total cash for admins
 */

UPDATE
  user_total ut
  INNER JOIN
    (SELECT
      t.`admin_allotted_id` AS user_allotted_id,
      SUM(t.`amount`) AS amount
    FROM
      `transaction` t
    GROUP BY t.`admin_allotted_id`) a
    ON 1 = 1
    AND ut.`user_id` = a.user_allotted_id 
    SET ut.`amount` = a.amount
WHERE ut.`user_id` = a.user_allotted_id;


/**
 * Repopulate user total coins coming from coins given
 */

UPDATE
  user_total ut
  INNER JOIN
    (SELECT
      ct.`user_id`,
      SUM(ABS(amount)) AS amount
    FROM
      coin_transaction ct
    WHERE ct.`transaction_id` IS NULL
    GROUP BY ct.`user_id`) a
    	ON 1 = 1
    		AND ut.`user_id` = a.user_id 
    SET ut.`coin_amount` = a.amount;

/**
 * Repopulate coins coming from orders 
 */

-- Analyze orders, insert coins related to orders
-- coin transactions for admin, user, and seller
SET @admin_id = (SELECT
    MAX(u.id)
  FROM
    `user` u
    INNER JOIN user_type ut
      ON 1 = 1
      AND u.user_type_id = ut.id
      AND ut.name = 'Admin');
      
SET @PACKAGE_PURCHASE_DEDUCTION = (SELECT CAST(VALUE AS DECIMAL(65,2)) FROM sysparam WHERE `key` = 'PACKAGE_PURCHASE_DEDUCTION');

INSERT INTO coin_transaction(created_by, created_date, is_active, order_id, `type`, amount, user_id)
-- dropshipper
SELECT 
  @admin_id AS created_by,
  o.created_date,
  o.is_active,
  o.id AS order_id,
  'C' AS `type`,
  -75,
  o.dropshipper_id AS user_id
FROM `order` o 
WHERE 1 = 1

UNION ALL

-- seller
SELECT 
  @admin_id AS created_by,
  o.created_date,
  o.is_active,
  o.id AS order_id,
  'D' AS `type`,
  ABS(@PACKAGE_PURCHASE_DEDUCTION),
  o.seller_id AS user_id
FROM `order` o 
WHERE 1 = 1

UNION ALL

-- admin
SELECT 
  @admin_id AS created_by,
  o.created_date,
  o.is_active,
  o.id AS order_id,  
  'D' AS `type`,
  ABS(@PACKAGE_PURCHASE_DEDUCTION) - 75,
  @admin_id AS user_id
FROM `order` o 
WHERE 1 = 1;

/**
 * Update user totals for orders
 */
UPDATE user_total ut 
INNER JOIN (SELECT 
  a.user_id,
  CASE 
    WHEN ut.name = 'Dropshipper' THEN ABS(SUM(a.amount))
    WHEN ut.name = 'Seller' THEN ABS(SUM(a.amount)) * -1
    WHEN ut.name = 'Admin' THEN ABS(SUM(a.amount))
  END AS amount,
  ut.`name` AS user_type
FROM (
SELECT 
  @admin_id AS created_by,
  o.created_date,
  o.is_active,
  o.id AS order_id,
  'C' AS `type`,
  -75 AS amount,
  o.dropshipper_id AS user_id
FROM `order` o 
WHERE 1 = 1

UNION ALL

-- seller
SELECT 
  @admin_id AS created_by,
  o.created_date,
  o.is_active,
  o.id AS order_id,
  'D' AS `type`,
  ABS(@PACKAGE_PURCHASE_DEDUCTION)  AS amount,
  o.seller_id AS user_id
FROM `order` o 
WHERE 1 = 1

UNION ALL

-- admin
SELECT 
  @admin_id AS created_by,
  o.created_date,
  o.is_active,
  o.id AS order_id,  
  'D' AS `type`,
  ABS(@PACKAGE_PURCHASE_DEDUCTION) - 75 AS amount,
  @admin_id AS user_id
FROM `order` o 
WHERE 1 = 1
) AS a


INNER JOIN `user` u 
  ON 1 = 1
    AND a.user_id = u.id
 INNER JOIN user_type ut
   ON 1 =1
     AND u.user_type_id = ut.`id`

GROUP BY user_id) a
  ON ut.`user_id` = a.user_id
  SET coin_amount = coin_amount + a.amount;
  
/**
 * Insert coin transactions for deliveries
 * 1. Rejected:  
 * 2. Add coin deductions for non delivered
 * 3. Add coin 
 */

-- Rejected: (Dropship & Parcel)

SET @admin_id = (SELECT
    MAX(u.id)
  FROM
    `user` u
    INNER JOIN user_type ut
      ON 1 = 1
      AND u.user_type_id = ut.id
      AND ut.name = 'Admin');
      
INSERT INTO coin_transaction(created_by, created_date, is_active, delivery_id, `type`, amount, user_id)
SELECT
  @admin_id AS created_by,  
  dt.`last_updated` AS created_date,
  0 AS is_active,
  d.id AS delivery_id,
  CASE 
    WHEN ut.name = 'Admin' THEN 'D'
    WHEN ut.name = 'Dropshipper' THEN 'C'
    WHEN ut.name = 'Seller' THEN 'D'
  END AS `type`,    
  CASE 
    WHEN ut.name = 'Admin' THEN d.service_fee - 35
    WHEN ut.name = 'Dropshipper' THEN -35
    WHEN ut.name = 'Seller' THEN d.service_fee
  END AS `amount`,
  CASE 
    WHEN ut.name = 'Admin' THEN @admin_id
    WHEN ut.name = 'Dropshipper' THEN d.dropshipper_id
    WHEN ut.name = 'Seller' THEN d.seller_id
  END AS `user_id`
FROM
  delivery d
  INNER JOIN delivery_status ds
    ON 1 = 1
    AND ds.`name` = 'Rejected'
    AND d.`delivery_status_id` = ds.`id`
  INNER JOIN delivery_tracking dt
    ON 1 = 1
    AND dt.`delivery_id` = d.`id`
  INNER JOIN delivery_status ds2
    ON 1 = 1
    AND ds2.`id` = dt.`delivery_status_id`
    AND ds2.`name` IN ('Rejected')
  INNER JOIN delivery_option `do`
    ON 1 = 1
      AND d.`delivery_option_id` = do.`id`
   CROSS JOIN user_type ut
ORDER BY d.`id` ASC, ds2.name, dt.last_updated, ut.`name`;

-- Non Rejected: (Dropship & Parcel)
INSERT INTO coin_transaction(created_by, created_date, is_active, delivery_id, `type`, amount, user_id)
SELECT
  @admin_id AS created_by,  
  d.created_date AS created_date,
  1 AS is_active,
  d.id AS delivery_id,
  CASE 
    WHEN ut.name = 'Admin' THEN 'D'
    WHEN ut.name = 'Dropshipper' THEN 'C'
    WHEN ut.name = 'Seller' THEN 'D'
  END AS `type`,    
  CASE 
    WHEN ut.name = 'Admin' THEN d.service_fee - 35
    WHEN ut.name = 'Dropshipper' THEN -35
    WHEN ut.name = 'Seller' THEN d.service_fee
  END AS `amount`,
  CASE 
    WHEN ut.name = 'Admin' THEN @admin_id
    WHEN ut.name = 'Dropshipper' THEN d.dropshipper_id
    WHEN ut.name = 'Seller' THEN d.seller_id
  END AS `user_id`
FROM
  delivery d
  INNER JOIN delivery_status ds
    ON 1 = 1
    AND ds.`name` NOT IN ('Rejected')
    AND d.`delivery_status_id` = ds.`id`
  INNER JOIN delivery_option `do`
    ON 1 = 1
      AND d.`delivery_option_id` = do.`id`
   CROSS JOIN user_type ut

ORDER BY d.`id` ASC, ds.name, d.created_date, ut.`name`;

-- Update totals 
UPDATE user_total ut 
INNER JOIN (
SELECT 
  user_id,
  CASE 
    WHEN user_type = 'Admin' THEN amount
    WHEN user_type = 'Dropshipper' THEN amount
    WHEN user_type = 'Seller' THEN amount * -1
  END AS amount
FROM (SELECT 
  a.user_id,
  a.user_type,
  ABS(SUM(a.amount)) AS amount
FROM (SELECT
  @admin_id AS created_by,  
  d.created_date AS created_date,
  1 AS is_active,
  d.id AS delivery_id,
  ut.name AS user_type,
  CASE 
    WHEN ut.name = 'Admin' THEN 'D'
    WHEN ut.name = 'Dropshipper' THEN 'C'
    WHEN ut.name = 'Seller' THEN 'D'
  END AS `type`,    
  CASE 
    WHEN ut.name = 'Admin' THEN d.service_fee - 35
    WHEN ut.name = 'Dropshipper' THEN -35
    WHEN ut.name = 'Seller' THEN d.service_fee
  END AS `amount`,
  CASE 
    WHEN ut.name = 'Admin' THEN @admin_id
    WHEN ut.name = 'Dropshipper' THEN d.dropshipper_id
    WHEN ut.name = 'Seller' THEN d.seller_id
  END AS `user_id`
FROM
  delivery d
  INNER JOIN delivery_status ds
    ON 1 = 1
    AND ds.`name` != 'Rejected'
    AND d.`delivery_status_id` = ds.`id`
  INNER JOIN delivery_option `do`
    ON 1 = 1
      AND d.`delivery_option_id` = do.`id`
   CROSS JOIN user_type ut
  
ORDER BY d.`id` ASC, ds.name, d.created_date, ut.`name`) AS a
GROUP BY user_id) AS a) AS a 
ON 1 = 1
  AND ut.`user_id` = a.user_id
  SET ut.coin_amount = ut.coin_amount + a.amount;

/**
  Delivered items (dropship)
  1. Use diff logic for "Parcel" and Dropship
  2. Just copy those above and reiinsert them all!
 */
INSERT INTO coin_transaction(created_by, created_date, is_active, delivery_id, `type`, amount, user_id)
SELECT
  @admin_id AS created_by,
  d.created_date AS created_date,
  1 AS is_active,
  d.id AS delivery_id,
  CASE
    WHEN ut.name = 'Admin' THEN 'C'
    WHEN ut.name = 'Dropshipper' THEN 'C'
    WHEN ut.name = 'Seller' THEN 'C'
  END AS `type`,
  CASE
    WHEN ut.name = 'Admin' THEN ABS(d.declared_amount) * -1
    WHEN ut.name = 'Dropshipper' THEN ABS(d.amount_distributor) * -1
    WHEN ut.name = 'Seller' THEN ABS(d.declared_amount - d.amount_distributor) * -1
  END AS `amount`,
  CASE
    WHEN ut.name = 'Admin' THEN @admin_id
    WHEN ut.name = 'Dropshipper' THEN d.dropshipper_id
    WHEN ut.name = 'Seller' THEN d.seller_id
  END AS `user_id`
FROM
  delivery d
  INNER JOIN delivery_status ds
    ON 1 = 1
    AND ds.`name` IN ('Delivered')
    AND d.`delivery_status_id` = ds.`id`
  INNER JOIN delivery_option `do`
    ON 1 = 1
      AND d.`delivery_option_id` = do.`id`
      AND do.name = 'Dropship'
   CROSS JOIN user_type ut

ORDER BY d.`id` ASC, ds.name, d.created_date, ut.`name`;

-- Update totals
UPDATE user_total ut
INNER JOIN (
SELECT
  user_id,
  CASE
    WHEN user_type = 'Admin' THEN amount
    WHEN user_type = 'Dropshipper' THEN amount
    WHEN user_type = 'Seller' THEN amount * -1
  END AS amount
FROM (SELECT
  a.user_id,
  a.user_type,
  ABS(SUM(a.amount)) AS amount
FROM (SELECT
  @admin_id AS created_by,
  d.created_date AS created_date,
  1 AS is_active,
  d.id AS delivery_id,
  ut.name AS user_type,
  CASE
    WHEN ut.name = 'Admin' THEN 'D'
    WHEN ut.name = 'Dropshipper' THEN 'C'
    WHEN ut.name = 'Seller' THEN 'D'
  END AS `type`,
  CASE
    WHEN ut.name = 'Admin' THEN d.service_fee - 35
    WHEN ut.name = 'Dropshipper' THEN -35
    WHEN ut.name = 'Seller' THEN d.service_fee
  END AS `amount`,
  CASE
    WHEN ut.name = 'Admin' THEN @admin_id
    WHEN ut.name = 'Dropshipper' THEN d.dropshipper_id
    WHEN ut.name = 'Seller' THEN d.seller_id
  END AS `user_id`
FROM
  delivery d
  INNER JOIN delivery_status ds
    ON 1 = 1
    AND ds.`name` = 'Delivered'
    AND d.`delivery_status_id` = ds.`id`
  INNER JOIN delivery_option `do`
    ON 1 = 1
      AND d.`delivery_option_id` = do.`id`
      AND do.name = 'Dropship'
   CROSS JOIN user_type ut

ORDER BY d.`id` ASC, ds.name, d.created_date, ut.`name`) AS a
GROUP BY user_id) AS a) AS a
ON 1 = 1
  AND ut.`user_id` = a.user_id
  SET ut.coin_amount = ut.coin_amount + a.amount;

/**
  Delivered items (Parcel)
 */
INSERT INTO coin_transaction(created_by, created_date, is_active, delivery_id, `type`, amount, user_id)
SELECT
  @admin_id AS created_by,
  d.created_date AS created_date,
  1 AS is_active,
  d.id AS delivery_id,
  CASE
    WHEN ut.name = 'Admin' THEN 'C'
    WHEN ut.name = 'Seller' THEN 'C'
  END AS `type`,
  CASE
    WHEN ut.name = 'Admin' THEN ABS(d.declared_amount) * -1
    WHEN ut.name = 'Seller' THEN ABS(d.declared_amount) * -1
  END AS `amount`,
  CASE
    WHEN ut.name = 'Admin' THEN @admin_id
    WHEN ut.name = 'Seller' THEN d.seller_id
  END AS `user_id`
FROM
  delivery d
  INNER JOIN delivery_status ds
    ON 1 = 1
    AND ds.`name` IN ('Delivered')
    AND d.`delivery_status_id` = ds.`id`
  INNER JOIN delivery_option `do`
    ON 1 = 1
      AND d.`delivery_option_id` = do.`id`
      AND do.name = 'Parcel'
  CROSS JOIN user_type ut
    ON  1 = 1
      AND ut.name IN ('Admin', 'Seller')

ORDER BY d.`id` ASC, ds.name, d.created_date, ut.`name`;

-- Update totals
UPDATE user_total ut
INNER JOIN (
SELECT
  user_id,
  CASE
    WHEN user_type = 'Admin' THEN amount
    WHEN user_type = 'Dropshipper' THEN amount
    WHEN user_type = 'Seller' THEN amount * -1
  END AS amount
FROM (SELECT
  a.user_id,
  a.user_type,
  ABS(SUM(a.amount)) AS amount
FROM (SELECT
  @admin_id AS created_by,
  d.created_date AS created_date,
  1 AS is_active,
  d.id AS delivery_id,
  ut.name AS user_type,
  CASE
    WHEN ut.name = 'Admin' THEN 'D'
    WHEN ut.name = 'Dropshipper' THEN 'C'
    WHEN ut.name = 'Seller' THEN 'D'
  END AS `type`,
  CASE
    WHEN ut.name = 'Admin' THEN d.service_fee - 35
    WHEN ut.name = 'Dropshipper' THEN -35
    WHEN ut.name = 'Seller' THEN d.service_fee
  END AS `amount`,
  CASE
    WHEN ut.name = 'Admin' THEN @admin_id
    WHEN ut.name = 'Dropshipper' THEN d.dropshipper_id
    WHEN ut.name = 'Seller' THEN d.seller_id
  END AS `user_id`
FROM
  delivery d
  INNER JOIN delivery_status ds
    ON 1 = 1
    AND ds.`name` = 'Delivered'
    AND d.`delivery_status_id` = ds.`id`
  INNER JOIN delivery_option `do`
    ON 1 = 1
      AND d.`delivery_option_id` = do.`id`
      AND do.name = 'Parcel'
   CROSS JOIN user_type ut
     ON  1 = 1
       AND ut.name IN ('Admin', 'Seller')

ORDER BY d.`id` ASC, ds.name, d.created_date, ut.`name`) AS a
GROUP BY user_id) AS a) AS a
ON 1 = 1
  AND ut.`user_id` = a.user_id
  SET ut.coin_amount = ut.coin_amount + a.amount;

/**
  Final update totals
 */

-- Form updates here
SELECT
    user_id,
    ut.`name` AS user_type,
    CASE
        WHEN ut.name = 'Admin' THEN SUM(ct.amount)
        WHEN ut.name = 'Dropshipper' THEN SUM(ct.amount) * -1
        WHEN ut.name = 'Seller' THEN SUM(ct.amount) * -1
        END AS `amount`
FROM
    coin_transaction ct
        INNER JOIN `user` u
                   ON 1= 1
                       AND ct.`user_id` = u.`id`
        INNER JOIN user_type ut
                   ON 1 =1
                       AND u.`user_type_id` = ut.`id`
WHERE 1 = 1
  AND ct.is_active = 1
GROUP BY u.`id`