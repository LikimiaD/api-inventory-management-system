SELECT order_id, status, created_at
FROM orders
WHERE customer_id = $1
LIMIT COALESCE($2, 1000);