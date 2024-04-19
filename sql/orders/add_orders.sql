INSERT INTO orders (customer_id, status, created_at, updated_at)
VALUES ($1, 'new', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING order_id;