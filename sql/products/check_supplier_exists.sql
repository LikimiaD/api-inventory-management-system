SELECT EXISTS (
    SELECT 1
    FROM suppliers
    WHERE supplier_id = $1
);