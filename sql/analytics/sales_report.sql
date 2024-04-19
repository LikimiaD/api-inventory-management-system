SELECT
    p.product_id,
    p.name,
    COALESCE(SUM(od.price * od.quantity), 0) AS total_sales
FROM
    products p
        LEFT JOIN
    order_details od ON p.product_id = od.product_id
GROUP BY
    p.product_id, p.name;