UPDATE products
SET
    name = COALESCE($2, name),
    supplier_id = COALESCE($3, supplier_id),
    description = COALESCE($4, description),
    price = COALESCE($5, price),
    quantity = COALESCE($6, quantity),
    category = COALESCE($7, category),
    updated_at = CURRENT_TIMESTAMP
WHERE product_id =   $1;