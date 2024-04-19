SELECT product_id, avg(quantity) as avg_sold
FROM order_details
GROUP BY product_id;