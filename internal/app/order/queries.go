package order

const getOrderByTableQuery = `
SELECT
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
FROM ceng_order o
JOIN ceng_table t ON o.table_id = t.id
JOIN ceng_course c ON c.order_id = o.id
JOIN ceng_course_selection cs ON cs.course_id = c.id AND cs.quantity > 0
JOIN ceng_menu_item mi ON mi.id = cs.menu_item_id
JOIN ceng_menu_category mc ON mc.id = mi.menu_category_id
JOIN ceng_printer p ON p.id = mc.printer_id
LEFT JOIN ceng_menu_option mo ON mo.id = cs.menu_option_id
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
FROM ceng_order o
JOIN ceng_table t ON o.table_id = t.id
JOIN ceng_course c ON c.order_id = o.id
JOIN ceng_course_selection cs ON cs.course_id = c.id AND cs.quantity > 0
JOIN ceng_menu_item mi ON mi.id = cs.menu_item_id
JOIN ceng_menu_category mc ON mc.id = mi.menu_category_id
JOIN ceng_printer p ON p.id = mc.printer_id
LEFT JOIN ceng_menu_option mo ON mo.id = cs.menu_option_id
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
    t.name AS table_name,
    t.created_at AS table_created_at,
    p.id AS printer_id,
    'PRECONTO' AS printer_title,
    p.url AS printer_url,
    MIN(c.id) AS course_id,
    0 AS course_number, 
    mi.title AS menu_item_title,
    mi.price AS menu_item_price,
    mo.title AS menu_option_title,
    mo.price AS menu_option_price,
    SUM(cs.quantity) AS quantity
FROM ceng_order o
JOIN ceng_table t ON o.table_id = t.id
JOIN ceng_course c ON c.order_id = o.id
JOIN ceng_course_selection cs ON cs.course_id = c.id AND cs.quantity > 0
JOIN ceng_menu_item mi ON mi.id = cs.menu_item_id
JOIN ceng_menu_category mc ON mc.id = mi.menu_category_id
JOIN ceng_printer p ON (
    (t.inside = TRUE  AND p.title = 'PRECONTO')
    OR (t.inside = FALSE AND p.title = 'PRECONTO_ASPORTO')
)
LEFT JOIN ceng_menu_option mo ON mo.id = cs.menu_option_id
WHERE o.table_id = $1
GROUP BY
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
    t.name AS table_name,
    t.created_at AS table_created_at,
    t.payment_method as table_payment,
    p.id AS printer_id,
    'PAGAMENTO' AS printer_title,
    p.url AS printer_url,
    SUM(cs.quantity * COALESCE(mo.price, mi.price)) AS price_total
FROM ceng_order o
JOIN ceng_table t ON o.table_id = t.id
JOIN ceng_course c ON c.order_id = o.id
JOIN ceng_course_selection cs ON cs.course_id = c.id AND cs.quantity > 0
JOIN ceng_menu_item mi ON mi.id = cs.menu_item_id
JOIN ceng_menu_category mc ON mc.id = mi.menu_category_id
JOIN ceng_printer p ON (
    (t.inside = TRUE  AND p.title = 'PAGAMENTO')
    OR (t.inside = FALSE AND p.title = 'PAGAMENTO_ASPORTO')
)
LEFT JOIN ceng_menu_option mo ON mo.id = cs.menu_option_id
WHERE o.table_id = $1
AND p.active = true
GROUP BY
    t.name,
    t.created_at,
    t.payment_method,
    p.id,
    p.title,
    p.url
`
