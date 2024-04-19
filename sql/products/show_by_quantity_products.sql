SELECT product_id, supplier_id, name, description, price, quantity, category, created_at, updated_at FROM products
WHERE quantity > 0
LIMIT COALESCE($1, 1000);