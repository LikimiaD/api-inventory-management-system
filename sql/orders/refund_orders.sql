WITH updated_orders AS (
    UPDATE orders SET status = 'refunded' WHERE order_id = $1 RETURNING order_id
), updated_products AS (
    UPDATE products
        SET quantity = products.quantity + od.quantity
        FROM order_details AS od
        WHERE od.order_id = $1 AND products.product_id = od.product_id
        RETURNING products.product_id
)
SELECT 1;