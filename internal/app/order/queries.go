package order

const getOrderByTableQuery = `
SELECT
    u.username AS username,
    t.name AS table_name,
    t.created_at AS table_created_at,
    p.id AS printer_id,
    p.title AS printer_title,
    p.url AS printer_url,
    c.id AS course_id,
    c.position AS course_number,
    mi.title AS menu_item_title,
    mi.price AS menu_item_price,
    mo.title AS menu_option_title,
    mo.price AS menu_option_price,
    cs.quantity
FROM hueat_order o
JOIN hueat_table t ON o.table_id = t.id
JOIN hueat_user u ON t.user_id = u.id
JOIN hueat_course c ON c.order_id = o.id
JOIN hueat_course_selection cs ON cs.course_id = c.id AND cs.quantity > 0
JOIN hueat_menu_item mi ON mi.id = cs.menu_item_id
JOIN hueat_menu_category mc ON mc.id = mi.menu_category_id
JOIN hueat_printer p ON  (
    (t.inside = TRUE AND p.id = mi.printer_inside_id)
    OR (t.inside = FALSE AND p.id = mi.printer_outside_id)
)
LEFT JOIN hueat_menu_option mo ON mo.id = cs.menu_option_id
WHERE o.table_id = $1
AND p.active = true
ORDER BY
    p.id,       
    c.position, 
    mc.position,
    mi.position,
    mo.position;
`

const getCourseByTableAndCourseQuery = `
SELECT
    u.username AS username,
    t.name AS table_name,
    t.created_at AS table_created_at,
    p.id AS printer_id,
    p.title AS printer_title,
    p.url AS printer_url,
    c.id AS course_id,
    c.position AS course_number,
    mi.title AS menu_item_title,
    mi.price AS menu_item_price,
    mo.title AS menu_option_title,
    mo.price AS menu_option_price,
    cs.quantity
FROM hueat_order o
JOIN hueat_table t ON o.table_id = t.id
JOIN hueat_user u ON t.user_id = u.id
JOIN hueat_course c ON c.order_id = o.id
JOIN hueat_course_selection cs ON cs.course_id = c.id AND cs.quantity > 0
JOIN hueat_menu_item mi ON mi.id = cs.menu_item_id
JOIN hueat_menu_category mc ON mc.id = mi.menu_category_id
JOIN hueat_printer p ON  (
    (t.inside = TRUE AND p.id = mi.printer_inside_id)
    OR (t.inside = FALSE AND p.id = mi.printer_outside_id)
)
LEFT JOIN hueat_menu_option mo ON mo.id = cs.menu_option_id
WHERE o.table_id = $1
AND c.id = $2
AND p.active = true
ORDER BY
    p.id,       
    c.position, 
    mc.position,
    mi.position,
    mo.position;
`

const getPricedOrderByTableQuery = `
SELECT
    u.username AS username,
    t.name AS table_name,
    t.created_at AS table_created_at,
    p.id AS printer_id,
    (SELECT cfg.config_value FROM hueat_config cfg WHERE cfg.config_key = 'pricedOrderPrinterTitle') AS printer_title,
    p.url AS printer_url,
    MIN(c.id) AS course_id,
    0 AS course_number, 
    mi.title AS menu_item_title,
    mi.price AS menu_item_price,
    mo.title AS menu_option_title,
    mo.price AS menu_option_price,
    SUM(cs.quantity) AS quantity
FROM hueat_order o
JOIN hueat_table t ON o.table_id = t.id
JOIN hueat_user u ON t.user_id = u.id
JOIN hueat_course c ON c.order_id = o.id
JOIN hueat_course_selection cs ON cs.course_id = c.id AND cs.quantity > 0
JOIN hueat_menu_item mi ON mi.id = cs.menu_item_id
JOIN hueat_menu_category mc ON mc.id = mi.menu_category_id
JOIN hueat_config cfg ON (
    (t.inside = TRUE AND cfg.config_key = 'pricedOrderPrinterInsideID')
    OR (t.inside = FALSE AND cfg.config_key = 'pricedOrderPrinterOutsideID')
)
JOIN hueat_printer p ON p.id = cfg.config_value
LEFT JOIN hueat_menu_option mo ON mo.id = cs.menu_option_id
WHERE o.table_id = $1
GROUP BY
    u.username,
    t.name,
    t.created_at,
    p.id,
    p.title,
    p.url,
    mi.title,
    mi.price,
    mo.title,
    mo.price,
    mi.position,
    mc.position,
    mo.position
ORDER BY
    mc.position,
    mi.position,
    mo.position;
`

const getTotalPriceAndPaymentByTableQuery = `
SELECT
    u.username AS username,
    t.name AS table_name,
    t.created_at AS table_created_at,
    t.payment_method as table_payment,
    p.id AS printer_id,
    (SELECT cfg.config_value FROM hueat_config cfg WHERE cfg.config_key = 'totalPricePaymentPrinterTitle') AS printer_title,
    p.url AS printer_url,
    SUM(cs.quantity * COALESCE(mo.price, mi.price)) AS price_total
FROM hueat_order o
JOIN hueat_table t ON o.table_id = t.id
JOIN hueat_user u ON t.user_id = u.id
JOIN hueat_course c ON c.order_id = o.id
JOIN hueat_course_selection cs ON cs.course_id = c.id AND cs.quantity > 0
JOIN hueat_menu_item mi ON mi.id = cs.menu_item_id
JOIN hueat_config cfg ON (
    (t.inside = TRUE AND cfg.config_key = 'totalPricePaymentPrinterInsideID')
    OR (t.inside = FALSE AND cfg.config_key = 'totalPricePaymentPrinterOutsideID')
)
JOIN hueat_printer p ON p.id = cfg.config_value
LEFT JOIN hueat_menu_option mo ON mo.id = cs.menu_option_id
WHERE o.table_id = $1
AND p.active = true
GROUP BY
    u.username,
    t.name,
    t.created_at,
    t.payment_method,
    p.id,
    p.title,
    p.url
`
