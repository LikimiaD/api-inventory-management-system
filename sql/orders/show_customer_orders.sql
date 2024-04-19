SELECT order_id, customer_id, status, created_at, updated_at FROM orders
WHERE customer_id = $1
LIMIT COALESCE($2, 1000);