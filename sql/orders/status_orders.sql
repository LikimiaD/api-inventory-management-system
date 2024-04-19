UPDATE orders
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE order_id = $1;