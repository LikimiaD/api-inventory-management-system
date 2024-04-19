SELECT product_id, supplier_id, name, description, price, quantity, category, created_at, updated_at FROM products
WHERE price BETWEEN $1 AND $2
LIMIT COALESCE($3, 1000)