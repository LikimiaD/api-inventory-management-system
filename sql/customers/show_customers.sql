SELECT customer_id, name, email, phone, address
FROM customers
LIMIT COALESCE($1, 1000);