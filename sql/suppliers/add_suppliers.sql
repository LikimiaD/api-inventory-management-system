INSERT INTO suppliers (name, contact_name, contact_email, contact_phone)
VALUES ($1, $2, $3, $4) RETURNING supplier_id;