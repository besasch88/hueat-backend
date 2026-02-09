package statistics

const getPaymentMethodsTakinsQuery = `
SELECT t.payment_method as payment_method, SUM(COALESCE(o.price, i.price)) as takings
FROM ceng_course_selection AS s
INNER JOIN ceng_course AS c
ON s.course_id = c.id
INNER JOIN ceng_order as ord
ON c.order_id = ord.id
INNER JOIN ceng_table as t
ON ord.table_id = t.id
INNER JOIN ceng_menu_item AS i
ON s.menu_item_id = i.id
LEFT JOIN ceng_menu_option AS o
ON s.menu_option_id = o.id
WHERE t."close" IS TRUE
GROUP BY t.payment_method
ORDER BY takings DESC;
`

const getAverageTableDurationQuery = `
SELECT (EXTRACT(EPOCH FROM AVG(t.updated_at - t.created_at)) * 1e9)::bigint as avg_duration
FROM ceng_table as t
WHERE t.inside = TRUE
GROUP BY t."close"
HAVING t."close" is TRUE;
`

const getMenuItemStatsQuery = `
SELECT 
	COALESCE(o.title, i.title) as title, 
	SUM(s.quantity) as quantity, 
	SUM(s.quantity) * COALESCE(o.price, i.price) as takings
FROM ceng_course_selection AS s
INNER JOIN ceng_course AS c
ON s.course_id = c.id
INNER JOIN ceng_order as ord
ON c.order_id = ord.id
INNER JOIN ceng_table as t
ON ord.table_id = t.id
INNER JOIN ceng_menu_item AS i
ON s.menu_item_id = i.id
LEFT JOIN ceng_menu_option AS o
ON s.menu_option_id = o.id
WHERE t."close" IS TRUE
GROUP BY i.title, o.title, i.price, o.price
ORDER BY quantity DESC;
`
