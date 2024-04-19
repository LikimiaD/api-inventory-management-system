SELECT supplier_id, name, contact_name, contact_email, contact_phone
FROM suppliers
LIMIT COALESCE($1, 1000);