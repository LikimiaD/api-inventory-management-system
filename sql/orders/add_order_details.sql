INSERT INTO order_details (order_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4) RETURNING order_detail_id;