SELECT EXISTS (
    SELECT 1
    FROM products
    WHERE name = $1
);