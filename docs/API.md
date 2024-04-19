
## Login

### POST /login

This endpoint authenticates a user and returns an access token required for operations that need authentication.

#### Request

**URL:** `/login`

**Method:** `POST`

**Content-Type:** `application/json`

**Body:**
```json
{
  "login": "user_login",
  "password": "user_password"
}
```

**Fields:**
- `login` (string): The user's login username.
- `password` (string): The user's password.

#### Response

**Content-Type:** `application/json`

- **Success Response:**

  **Code:** `200 OK`

  **Content:**
  ```json
  {
    "status": 200,
    "token": "eyJhbGciOiJIU..."
  }
  ```

- **Error Responses:**

  **Code:** `400 Bad Request`

  **Content:** Occurs if the request is missing login or password data, if the user does not exist, or if the password is incorrect. Examples:

  ```json
  {
    "error": "Login or password is empty"
  }
  ```

  or

  ```json
  {
    "error": "User does not exist"
  }
  ```

  or

  ```json
  {
    "error": "Incorrect password"
  }
  ```

  **Code:** `500 Internal Server Error`

  **Content:** Server-side error during token generation or other internal operations. Example:

  ```json
  {
    "error": "Internal server error"
  }
  ```

#### Description

- The user sends their login and password.
- The server verifies the provided data, checks if the user exists and if the password is correct.
- If the validations are successful, the server generates a JWT (JSON Web Token) and returns it to the user.
- If there are errors during the process, the server returns an appropriate error code and description.

## Supplier

### 1. Add Supplier

**Endpoint:** `POST /add_supplier`

#### Authorization
- Requires a valid JWT obtained from the login process.

#### Request
**Content-Type:** `application/json`
**Body:**
```json
{
    "name": "test",
    "contact_name": "test",
    "contact_email": "test@example.com",
    "contact_phone": "800 555 35 35"
}
```

#### Response

**Success Response:**
- **Code:** `201 Created`
- **Content:**
```json
{
    "id": 2,
    "status": 201
}
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Not enough information to create"}`
- **Content:** `{"error": "There is a user with this email"}`
- **Content:** `{"error": "There is a user with this number"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Problem on the server side, please report the error time"}`

### 2. Delete Supplier

**Endpoint:** `POST /delete_supplier`

#### Authorization
- Requires a valid JWT.

#### Request
**Content-Type:** `application/json`
**Body:**
```json
{
    "id": 3
}
```

#### Response

**Success Response:**
- **Code:** `204 No Content`

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "No supplier found with the provided ID"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Problem on the server side, please report the error time"}`

### 3. Update Supplier

**Endpoint:** `POST /update_supplier`

#### Authorization
- Requires a valid JWT.

#### Request
**Content-Type:** `application/json`
**Body:**
```json
{
    "id": 2,
    "name": "Garry",
    "contact_name": null,
    "contact_email": null,
    "contact_phone": null
}
```

#### Response

**Success Response:**
- **Code:** `204 No Content`

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "No supplier found with the provided ID"}`
- **Content:** `{"error": "Invalid data format"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Problem on the server side, please report the error time"}`

### 4. Show Suppliers

**Endpoint:** `GET /show_suppliers`

#### Authorization
- Requires a valid JWT.

#### Query Parameters
- **format**: Specifies the output format (`json`, `csv`, `excel`).
- **limit**: Specifies the maximum number of suppliers to return.

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** `application/json`
- **Content-Type:** `text/csv`
- **Content-Type:** `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Invalid limit value"}`
- **Content:** `{"error": "Invalid format specified"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Failed to retrieve suppliers"}`
- **Content:** `{"error": "Error encoding response data"}`
- **Content:** `{"error": "Failed to generate CSV"}`
- **Content:** `{"error": "Failed to generate Excel file"}`

## Customer

### 1. Add Customer

**Endpoint:** `POST /add_customer`

#### Authorization
- Requires a valid JWT obtained from the login process.

#### Request
**Content-Type:** `application/json`
**Body:**
```json
{
    "name": "likimiad",
    "email": "likimiad@example.com",
    "phone": "800 555 35 35",
    "address": "Moscow, NITU MISIS"
}
```

#### Response

**Success Response:**
- **Code:** `201 Created`
- **Content:**
```json
{
    "id": 1,
    "status": 201
}
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** Various error messages such as "Not enough information to create", "There is a user with this email", "There is a user with this number"
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Problem on the server side, please report the error time"}`

### 2. Delete Customer

**Endpoint:** `POST /delete_customer`

#### Authorization
- Requires a valid JWT.

#### Request
**Content-Type:** `application/json`
**Body:**
```json
{
    "id": 3
}
```

#### Response

**Success Response:**
- **Code:** `204 No Content`

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "No customer found with the provided ID"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Problem on the server side, please report the error time"}`

### 3. Update Customer

**Endpoint:** `POST /update_customer`

#### Authorization
- Requires a valid JWT.

#### Request
**Content-Type:** `application/json`
**Body:**
```json
{
    "id": 2,
    "name": "Garry",
    "email": null,
    "phone": null,
    "address": null
}
```

#### Response

**Success Response:**
- **Code:** `204 No Content`

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "No customer found with the provided ID"}`
- **Content:** `{"error": "Invalid data format"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Problem on the server side, please report the error time"}`

### 4. Show Customers

**Endpoint:** `GET /show_customers`

#### Authorization
- Requires a valid JWT.

#### Query Parameters
- **format:** Specifies the output format (`json`, `csv`, `excel`).
- **limit:** Specifies the maximum number of customers to return.

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** Varies based on the format parameter (`application/json`, `text/csv`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)

**Content:**
```json
[
    {
        "customer_id": 1,
        "customer_name": "likimiad",
        "customer_email": "likimiad@example.com",
        "customer_phone": "800 555 35 35",
        "customer_address": "Moscow, NITU MISIS"
    },
    {
        "customer_id": 3,
        "customer_name": "test2",
        "customer_email": "test2@example.com",
        "customer_phone": "802 555 35 35",
        "customer_address": "Moscow, NITU MISIS"
    },
    {
        "customer_id": 2,
        "customer_name": "Garry",
        "customer_email": "test@example.com",
        "customer_phone": "801 555 35 35",
        "customer_address": "Moscow, NITU MISIS"
    }
]
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** Various error messages such as "Invalid limit value", "Invalid format specified"
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** Various error messages such as "Failed to retrieve customers", "Error encoding response data", "Failed to generate CSV", "Failed to generate Excel file"

## Product

### 1. Add Product

**Endpoint:** `POST /add_product`

#### Authorization
- Requires a valid JWT obtained from the login process.

#### Request
**Content-Type:** `application/json`
**Body:**
```json
{
    "supplier_id": 1,
    "name": "Tomato",
    "description": "Organic tomatoes from Spain",
    "price": 12.50,
    "quantity": 200,
    "category": "vegetable"
}
```

#### Response

**Success Response:**
- **Code:** `201 Created`
- **Content:**
```json
{
    "id": 1,
    "status": 201
}
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** Various error messages such as "Not enough information to create", "The price of the product cannot be negative", "Product quantities cannot be negative", "Product with name [product_name] exists", "Supplier with ID [supplier_id] does not exist"
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Problem on the server side, please report the error time"}`

### 2. Delete Product

**Endpoint:** `POST /delete_product`

#### Authorization
- Requires a valid JWT.

#### Request
**Content-Type:** `application/json`
**Body:**
```json
{
    "id": 5
}
```

#### Response

**Success Response:**
- **Code:** `204 No Content`

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "No product found with the provided ID"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Problem on the server side, please report the error time"}`

### 3. Update Product

**Endpoint:** `POST /update_product`

#### Authorization
- Requires a valid JWT.

#### Request
**Content-Type:** `application/json`
**Body:**
```json
{
    "id": 2,
    "supplier_id": 2,
    "name": "Banana",
    "description": null,
    "price": null,
    "quantity": null,
    "category": null
}
```

#### Response

**Success Response:**
- **Code:** `204 No Content`

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** Various error messages such as "Name cannot be empty", "Supplier ID cannot be empty", "The price of the product cannot be negative", "Product quantities cannot be negative", "Product with name [product_name] exists", "Supplier with ID [supplier_id] does not exist"
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Problem on the server side, please report the error time"}`

### 4. Show Products

**Endpoint:** `GET /show_products`

#### Authorization
- Requires a valid JWT.

#### Query Parameters
- **format:** Specifies the output format (`json`, `csv`, `excel`).
- **limit:** Specifies the maximum number of products to return.

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** Varies based on the format parameter (`application/json`, `text/csv`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)

**Content:** (example for `json`)
```json
[
    {
        "product_id": 1,
        "supplier_id": 1,
        "name": "Tomato",
        "description": "Organic tomatoes from Spain",
        "price": 12.5,
        "quantity": 200,
        "category": "vegetable",
        "created_at": "2024-04-19T10:55:27.470113Z",
        "updated_at": "2024-04-19T10:55:27.470113Z"
    }
]
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Invalid limit value"}`
- **Content:** `{"error": "Invalid format specified"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server

Error`
- **Content:** `{"error": "Failed to retrieve products"}`
- **Content:** `{"error": "Error encoding response data"}`
- **Content:** `{"error": "Failed to generate CSV"}`
- **Content:** `{"error": "Failed to generate Excel file"}`

### 5. Show Products In Stock

**Endpoint:** `GET /show_in_stock_products`

#### Authorization
- Requires a valid JWT.

#### Query Parameters
- **format:** Specifies the output format (`json`, `csv`, `excel`).
- **limit:** Specifies the maximum number of products to return.

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** Varies based on the format parameter (`application/json`, `text/csv`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)

**Content:** (example for `json`)
```json
[
    {
        "product_id": 1,
        "supplier_id": 1,
        "name": "Tomato",
        "description": "Organic tomatoes from Spain",
        "price": 12.5,
        "quantity": 200,
        "category": "vegetable",
        "created_at": "2024-04-19T10:55:27.470113Z",
        "updated_at": "2024-04-19T10:55:27.470113Z"
    }
]
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Invalid limit value"}`
- **Content:** `{"error": "Invalid format specified"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server

Error`
- **Content:** `{"error": "Failed to retrieve products"}`
- **Content:** `{"error": "Error encoding response data"}`
- **Content:** `{"error": "Failed to generate CSV"}`
- **Content:** `{"error": "Failed to generate Excel file"}`

### 6. Show Products By Category

**Endpoint:** `GET /show_category_products`

#### Authorization
- Requires a valid JWT.

#### Query Parameters
- **category:** Specifies the category to be displayed
- **format:** Specifies the output format (`json`, `csv`, `excel`).
- **limit:** Specifies the maximum number of products to return.

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** Varies based on the format parameter (`application/json`, `text/csv`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)

**Content:** (example for `json`)
```json
[
    {
        "product_id": 1,
        "supplier_id": 1,
        "name": "Tomato",
        "description": "Organic tomatoes from Spain",
        "price": 12.5,
        "quantity": 200,
        "category": "vegetable",
        "created_at": "2024-04-19T10:55:27.470113Z",
        "updated_at": "2024-04-19T10:55:27.470113Z"
    }
]
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Invalid limit value"}`
- **Content:** `{"error": "Invalid format specified"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server

Error`
- **Content:** `{"error": "Failed to retrieve products"}`
- **Content:** `{"error": "Error encoding response data"}`
- **Content:** `{"error": "Failed to generate CSV"}`
- **Content:** `{"error": "Failed to generate Excel file"}`

### 7. Show Products By Price Range

**Endpoint:** `GET /show_price_products`

#### Authorization
- Requires a valid JWT.

#### Query Parameters
- **min** Specifies the minimum value be displayed
- **max** Specifies the maximum value be displayed
- **format:** Specifies the output format (`json`, `csv`, `excel`).
- **limit:** Specifies the maximum number of products to return.

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** Varies based on the format parameter (`application/json`, `text/csv`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)

**Content:** (example for `json`)
```json
[
  {
    "product_id": 4,
    "supplier_id": 2,
    "name": "Potato1",
    "description": "Versatile potatoes from Ireland",
    "price": 5.2,
    "quantity": 90,
    "category": "vegetable",
    "created_at": "2024-04-19T10:56:00.897553Z",
    "updated_at": "2024-04-19T10:56:00.897553Z"
  },
  {
    "product_id": 7,
    "supplier_id": 1,
    "name": "Onion2",
    "description": "Sweet onions from Vidalia",
    "price": 6.75,
    "quantity": 130,
    "category": "vegetable",
    "created_at": "2024-04-19T10:56:14.332062Z",
    "updated_at": "2024-04-19T10:56:14.332062Z"
  }
]
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Invalid limit value"}`
- **Content:** `{"error": "Invalid format specified"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server

Error`
- **Content:** `{"error": "Failed to retrieve products"}`
- **Content:** `{"error": "Error encoding response data"}`
- **Content:** `{"error": "Failed to generate CSV"}`
- **Content:** `{"error": "Failed to generate Excel file"}`

### 8. Show Products For Purchase Request

**Endpoint:** `GET /show_purchase_request`

#### Authorization
- Requires a valid JWT.

#### Query Parameters
- **maxQuantity** Specifies the filter for the number of output items (if lower, output)
- **format:** Specifies the output format (`json`, `csv`, `excel`).
- **limit:** Specifies the maximum number of products to return.

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** Varies based on the format parameter (`application/json`, `text/csv`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)

**Content:** (example for `json`)
```json
[
  {
    "product_id": 8,
    "name": "Garlic1",
    "supplier_id": 2,
    "contact_email": "test@example.com"
  },
  {
    "product_id": 4,
    "name": "Potato1",
    "supplier_id": 2,
    "contact_email": "test@example.com"
  }
]
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Invalid limit value"}`
- **Content:** `{"error": "Invalid format specified"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server

Error`
- **Content:** `{"error": "Failed to retrieve products"}`
- **Content:** `{"error": "Error encoding response data"}`
- **Content:** `{"error": "Failed to generate CSV"}`
- **Content:** `{"error": "Failed to generate Excel file"}`

## Order

### 1. Add Order

**Endpoint:** `POST /add_order`

#### Authorization
- Requires a valid JWT token for authentication.

#### Request Body
- **customer_id:** ID of the customer placing the order.
- **status:** Current status of the order (`created`, `processing`, etc.).
- **product_id:** ID of the product being ordered.
- **quantity:** Number of units of the product.
- **price:** Price per unit of the product.

```json
{
    "customer_id": 1,
    "status": "created",
    "product_id": 1,
    "quantity": 100,
    "price": 10.0
}
```

#### Response

**Success Responses:**
- **Code:** `201 Created`
- **Content-Type:** `application/json`
- **Content:**
```json
{
    "order_detail_id": 1,
    "order_id": 1,
    "status": 201
}
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Not enough information to create"}`
- **Content:** `{"error": "Price cannot be negative"}`
- **Content:** `{"error": "Quantities cannot be negative"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Problem when adding a new order"}`

### 2. Refund Order

**Endpoint:** `POST /refund_order`

#### Authorization
- Requires a valid JWT token for authentication.

#### Request Body
- **order_id:** ID of the order to be refunded.

```json
{
    "order_id": 2
}
```

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** `application/json`
- **Content:**
```json
{
    "message": "Refund processed successfully for order 2",
    "status": 200
}
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Order ID is required"}`
- **Content:** `{"error": "Order with ID %d does not exist"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Failed to process refund for order %d"}`

### 3. Update Order Status

**Endpoint:** `POST /update_order`

#### Authorization
- Requires a valid JWT token for authentication.

#### Request Body
- **order_id:** ID of the order whose status is to be updated.
- **status:** New status of the order (`accepted`, `shipped`, etc.).

```json
{
    "order_id": 1,
    "status": "accepted"
}
```

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** `application/json`
- **Content:**
```json
{
    "message": "Order status updated successfully for order 1",
    "status": 200
}
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Order ID and Status are required"}`
- **Content:** `{"error": "Updating status from 'refund' is not safe and not allowed"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Failed to update order status"}`

### 4. Show Customer Orders

**Endpoint:** `GET /show_customer_orders`

#### Authorization
- Requires a valid JWT token for authentication.

#### Query Parameters
- **customer_id:** ID of the customer whose orders are to be shown.
- **format:** Specifies the output format (`json`, `csv`, `excel`).
- **limit:** Specifies the maximum number of orders to return.

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** Varies based on the format parameter (`application/json`, `text/csv`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
- **Content:** (example for `json`)
```json
[
    {
        "order_id": 3,
        "status": "new",
        "created_at": "2024-04-19T11:19:42.140541Z"
    },
    {
        "order_id": 2,
        "status": "refunded",
        "created_at": "2024-04-19T11:19:39.644021Z"
    },
    {
        "order_id": 1,


        "status": "accepted",
        "created_at": "2024-04-19T11:19:35.338317Z"
    }
]
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Customer ID is required"}`
- **Content:** `{"error": "Invalid customer ID format"}`
- **Content:** `{"error": "Invalid limit value"}`
- **Content:** `{"error": "Invalid format specified"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Failed to retrieve orders"}`
- **Content:** `{"error": "Error encoding response data"}`
- **Content:** `{"error": "Failed to generate CSV"}`
- **Content:** `{"error": "Failed to generate Excel file"}`

### 5. Show Customer Orders Full

**Endpoint:** `GET /show_customer_orders_full`

#### Authorization
- Requires a valid JWT token for authentication.

#### Query Parameters
- **customer_id:** ID of the customer whose full order details are to be shown.
- **format:** Specifies the output format (`json`, `csv`, `excel`).
- **limit:** Specifies the maximum number of orders to return.

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** Varies based on the format parameter (`application/json`, `text/csv`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
- **Content:** (example for `json`)
```json
[
    {
        "order_id": 3,
        "customer_id": 1,
        "status": "new",
        "created_at": "2024-04-19T11:19:42.140541Z",
        "updated_at": "2024-04-19T11:19:42.140541Z"
    },
    {
        "order_id": 2,
        "customer_id": 1,
        "status": "refunded",
        "created_at": "2024-04-19T11:19:39.644021Z",
        "updated_at": "2024-04-19T11:19:39.644021Z"
    },
    {
        "order_id": 1,
        "customer_id": 1,
        "status": "accepted",
        "created_at": "2024-04-19T11:19:35.338317Z",
        "updated_at": "2024-04-19T11:22:00.567468Z"
    }
]
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Customer ID is required"}`
- **Content:** `{"error": "Invalid customer ID format"}`
- **Content:** `{"error": "Invalid limit value"}`
- **Content:** `{"error": "Invalid format specified"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Failed to retrieve orders"}`

### 6. Show Orders by Date

**Endpoint:** `GET /show_orders_by_date`

#### Authorization
- Requires a valid JWT token for authentication.

#### Query Parameters
- **startDate:** Start date for the order search (ISO8601 format).
- **endDate:** End date for the order search (ISO8601 format).
- **limit:** Specifies the maximum number of orders to return.

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** `application/json`
- **Content:**
```json
[
    {
        "order_id": 3,
        "status": "new",
        "created_at": "2024-04-19T11:19:42.140541Z"
    },
    {
        "order_id": 4,
        "status": "new",
        "created_at": "2024-04-19T11:19:52.380921Z"
    },
    {
        "order_id": 5,
        "status": "new",
        "created_at": "2024-04-19T11:19:55.176833Z"
    },
    {
        "order_id": 6,
        "status": "new",
        "created_at": "2024-04-19T11:19:57.952829Z"
    },
    {
        "order_id": 2,
        "status": "refunded",
        "created_at": "2024-04-19T11:19:39.644021Z"
    },
    {
        "order_id": 1,
        "status": "accepted",
        "created_at": "2024-04-19T11:19:35.338317Z"
    }
]
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Both startDate and endDate parameters are required"}`
- **Content:** `{"error": "Invalid startDate format, use ISO8601 format"}`
- **Content:** `{"error": "Invalid endDate format, use ISO8601 format"}`
- **Content:** `{"error": "Invalid limit value"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Failed to retrieve orders"}`

### 7. Show Orders by Status

**Endpoint:** `GET /show_orders_by_status`

#### Authorization
- Requires a

valid JWT token for authentication.

#### Query Parameters
- **status:** Status of the orders to filter by.
- **limit:** Specifies the maximum number of orders to return.

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** `application/json`
- **Content:**
```json
[
    {
        "order_id": 2,
        "status": "refunded",
        "created_at": "2024-04-19T11:19:39.644021Z"
    }
]
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Status parameter is required"}`
- **Content:** `{"error": "Invalid limit value"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error: "Failed to retrieve orders"}`

### 8. Show Order Details API

**Endpoint:** `GET /show_order_details`

#### Authorization
- Requires a valid JWT token for authentication.

#### Query Parameters
- **order_id:** ID of the order for which details are requested.
- **format:** Specifies the output format (`json`, `csv`, `excel`).
- **limit:** Specifies the maximum number of detail entries to return.

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** Varies based on the format parameter (`application/json`, `text/csv`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
- **Content:** (example for `json`)
```json
[
    {
        "order_detail_id": 1,
        "quantity": 100,
        "price": 10,
        "name": "Tomato"
    }
]
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Order ID is required"}`
- **Content:** `{"error": "Invalid order ID format"}`
- **Content:** `{"error": "Invalid limit value"}`
- **Content:** `{"error": "Invalid format specified"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error: "Failed to retrieve order details"}`

## Analytics

### 1. Sales Report

**Endpoint:** `GET /sales_report`

#### Authorization
- Requires a valid JWT token for authentication.

#### Query Parameters
- **format:** Specifies the output format (`json`, `csv`, `excel`).

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** Varies based on the format parameter (`application/json`, `text/csv`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
- **Content:** (example for `json`)
```json
[
    {
        "product_id": 3,
        "avg_sold": 100
    },
    {
        "product_id": 6,
        "avg_sold": 100
    },
    {
        "product_id": 2,
        "avg_sold": 100
    },
    {
        "product_id": 7,
        "avg_sold": 100
    },
    {
        "product_id": 1,
        "avg_sold": 100
    },
    {
        "product_id": 8,
        "avg_sold": 100
    }
]
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Invalid format specified"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Failed to retrieve sales report"}`

### 2. Requirements Report

**Endpoint:** `GET /requirements_report`

#### Authorization
- Requires a valid JWT token for authentication.

#### Query Parameters
- **format:** Specifies the output format (`json`, `csv`, `excel`).

#### Response

**Success Responses:**
- **Code:** `200 OK`
- **Content-Type:** Varies based on the format parameter (`application/json`, `text/csv`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
- **Content:** (example for `json`)
```json
[
    {
        "product_id": 4,
        "name": "Potato1",
        "total_sales": 0
    },
    {
        "product_id": 6,
        "name": "Broccoli3",
        "total_sales": 0
    },
    {
        "product_id": 2,
        "name": "Banana",
        "total_sales": 0
    },
    {
        "product_id": 7,
        "name": "Onion2",
        "total_sales": 0
    },
    {
        "product_id": 3,
        "name": "Carrot2",
        "total_sales": 0
    },
    {
        "product_id": 1,
        "name": "Tomato",
        "total_sales": 0
    },
    {
        "product_id": 8,
        "name": "Garlic1",
        "total_sales": 0
    }
]
```

**Error Responses:**
- **Code:** `400 Bad Request`
- **Content:** `{"error": "Invalid format specified"}`
- **Code:** `401 Unauthorized`
- **Content:** `{"error": "Unauthorized"}`
- **Code:** `500 Internal Server Error`
- **Content:** `{"error": "Failed to retrieve requirements report"}`