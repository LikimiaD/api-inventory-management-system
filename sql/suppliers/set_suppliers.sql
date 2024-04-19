UPDATE suppliers
SET
    name = COALESCE($2, name),
    contact_name = COALESCE($3, contact_name),
    contact_email = COALESCE($4, contact_email),
    contact_phone = COALESCE($5, contact_phone)
WHERE supplier_id = $1;