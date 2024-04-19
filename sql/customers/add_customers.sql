INSERT INTO customers (name, email, phone, address)
VALUES ($1, $2, $3, $4) RETURNING customer_id;