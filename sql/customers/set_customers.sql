UPDATE customers
SET
    name = COALESCE($2, name),
    email = COALESCE($3, email),
    phone = COALESCE($4, phone),
    address = COALESCE($5, address)
WHERE customer_id = $1;