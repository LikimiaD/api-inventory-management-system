CREATE TABLE IF NOT EXISTS customers (
    customer_id INT GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(15) NOT NULL UNIQUE,
    address VARCHAR(255) NOT NULL,
    PRIMARY KEY(customer_id)
);
CREATE TABLE IF NOT EXISTS suppliers
(
    supplier_id   INT GENERATED ALWAYS AS IDENTITY,
    name          VARCHAR(255) NOT NULL,
    contact_name  VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255) NOT NULL UNIQUE,
    contact_phone VARCHAR(15)  NOT NULL UNIQUE,
    PRIMARY KEY (supplier_id)
);
CREATE TABLE IF NOT EXISTS products (
    product_id INT GENERATED ALWAYS AS IDENTITY,
    supplier_id INT NOT NULL,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    category VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(product_id),
    CONSTRAINT fk_suppliers
        FOREIGN KEY(supplier_id)
            REFERENCES suppliers(supplier_id)
);
CREATE TABLE IF NOT EXISTS orders (
    order_id INT GENERATED ALWAYS AS IDENTITY,
    customer_id INTEGER NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(order_id),
    CONSTRAINT fk_customer
        FOREIGN KEY(customer_id)
            REFERENCES customers(customer_id)
);
CREATE TABLE IF NOT EXISTS order_details (
    order_detail_id INT GENERATED ALWAYS AS IDENTITY,
    order_id INTEGER,
    product_id INTEGER,
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    PRIMARY KEY(order_detail_id),
    CONSTRAINT fk_orders
        FOREIGN KEY(order_id)
            REFERENCES orders(order_id),
    CONSTRAINT fk_products
        FOREIGN KEY(product_id)
            REFERENCES products(product_id)
);
CREATE TABLE IF NOT EXISTS trusted_users (
    trusted_user_id INT GENERATED ALWAYS AS IDENTITY,
    login VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_activity TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);