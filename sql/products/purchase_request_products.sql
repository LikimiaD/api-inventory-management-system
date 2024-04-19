SELECT p.product_id, p.name, s.supplier_id, s.contact_email
FROM products p
JOIN suppliers s ON p.supplier_id = s.supplier_id
WHERE p.quantity < $1;