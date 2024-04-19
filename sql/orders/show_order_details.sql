SELECT od.order_detail_id, od.quantity, od.price, p.name
FROM order_details od
JOIN products p ON od.product_id = p.product_id
WHERE od.order_id = $1
LIMIT COALESCE($2, 1000);