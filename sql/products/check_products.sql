SELECT product_id, supplier_id, name, description, price, quantity, category, created_at, updated_at
FROM products
WHERE quantity <= $1
LIMIT COALESCE($2, 1000);