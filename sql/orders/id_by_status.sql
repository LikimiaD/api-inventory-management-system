SELECT order_id, customer_id, created_at
FROM orders
WHERE status = $1
LIMIT COALESCE($2, 1000);