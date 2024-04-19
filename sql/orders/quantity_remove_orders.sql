UPDATE products
SET quantity = quantity - $2
WHERE product_id = $1;