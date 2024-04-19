INSERT INTO products (supplier_id, name, description, price, quantity, category, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING product_id;
