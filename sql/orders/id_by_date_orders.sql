SELECT order_id, status, created_at FROM orders
WHERE created_at BETWEEN $1 AND $2
LIMIT COALESCE($3, 1000);