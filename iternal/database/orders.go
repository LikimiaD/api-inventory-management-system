package database

import (
	"database/sql"
	"os"
	"time"
)

type OrderInfo struct {
	OrderID   int64     `json:"order_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	OrderID    int64     `json:"order_id"`
	CustomerID int64     `json:"customer_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderDetail struct {
	OrderDetailID int64   `json:"order_detail_id"`
	Quantity      int64   `json:"quantity"`
	Price         float64 `json:"price"`
	Name          string  `json:"name"`
}

func (db *Database) removeQuantity(productID, quantity int64) error {
	query, err := os.ReadFile(ordersPath + "quantity_remove_orders.sql")
	if err != nil {
		db.Log.Error("Database removeQuantity() -> Read SQL file", err)
		return err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database removeQuantity() -> db.Begin()", err)
		return err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	if _, err = tx.Exec(string(query), productID, quantity); err != nil {
		db.Log.Error("Database removeQuantity() -> tx.Exec()", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		db.Log.Error("Database ShowProducts() -> tx.Commit()", err)
		return err
	}
	committed = true

	return nil
}

func (db *Database) AddOrder(customerID, productID, quantity int64, price float64) (int64, int64, error) {
	queryOrder, err := os.ReadFile(ordersPath + "add_orders.sql")
	if err != nil {
		db.Log.Error("Database ShowProducts() -> Read SQL file", err)
		return 0, 0, err
	}
	queryOrderDetails, err := os.ReadFile(ordersPath + "add_order_details.sql")
	if err != nil {
		db.Log.Error("Database ShowProducts() -> Read SQL file", err)
		return 0, 0, err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database AddCustomer() -> db.Begin()", err)
		return 0, 0, err
	}

	var orderID, orderDetailID int64
	err = tx.QueryRow(string(queryOrder), customerID).Scan(&orderID)
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()
	if err != nil {
		db.Log.Error("Database AddOrder() -> QueryRow() order", err)
		return 0, 0, err
	}

	err = tx.QueryRow(string(queryOrderDetails), orderID, productID, quantity, price).Scan(&orderDetailID)
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()
	if err != nil {
		db.Log.Error("Database AddOrder() -> QueryRow() orderDetail", err)
		return 0, 0, err
	}

	if err = tx.Commit(); err != nil {
		db.Log.Error("Database AddOrder() -> tx.Commit()", err)
		return 0, 0, err
	}
	committed = true

	if err = db.removeQuantity(productID, quantity); err != nil {
		db.Log.Error("Database AddOrder() -> db.removeQuantity()", err)
		return 0, 0, err
	}

	return orderID, orderDetailID, nil
}

func (db *Database) readRowsOrderInfo(rows *sql.Rows) ([]OrderInfo, error) {
	defer func() {
		if err := rows.Close(); err != nil {
			db.Log.Error("Some troubles after rows.Close()", err)
		}
	}()

	var products []OrderInfo

	for rows.Next() {
		var o OrderInfo
		if err := rows.Scan(&o.OrderID, &o.Status, &o.CreatedAt); err != nil {
			db.Log.Error("Database readRowsOrderInfo() -> parsing rows", err)
			return nil, err
		}
		products = append(products, o)
	}

	if err := rows.Err(); err != nil {
		db.Log.Error("Database readRowsOrderInfo() -> rows.Err()", err)
		return nil, err
	}
	return products, nil
}

func (db *Database) ShowByCustomerOrders(customerID int64, limit *int) ([]OrderInfo, error) {
	query, err := os.ReadFile(ordersPath + "id_by_customer_orders.sql")
	if err != nil {
		db.Log.Error("Database ShowByCustomerOrders() -> Read SQL file", err)
		return nil, err
	}

	limitValue := sql.NullInt64{Valid: limit != nil}
	if limit != nil {
		limitValue.Int64 = int64(*limit)
	}

	rows, err := db.Query(string(query), customerID, limitValue)
	if err != nil {
		db.Log.Error("Database ShowByCustomerOrders() -> tx.Query()", err)
		return nil, err
	}

	return db.readRowsOrderInfo(rows)
}

func (db *Database) ShowByDateOrders(startDateRange, endDateRange time.Time, limit *int) ([]OrderInfo, error) {
	query, err := os.ReadFile(ordersPath + "id_by_date_orders.sql")
	if err != nil {
		db.Log.Error("Database ShowByDateOrders() -> Read SQL file", err)
		return nil, err
	}

	limitValue := sql.NullInt64{Valid: limit != nil}
	if limit != nil {
		limitValue.Int64 = int64(*limit)
	}

	rows, err := db.Query(string(query), startDateRange, endDateRange, limitValue)
	if err != nil {
		db.Log.Error("Database ShowByDateOrders() -> tx.Query()", err)
		return nil, err
	}

	return db.readRowsOrderInfo(rows)
}

func (db *Database) ShowByStatusOrders(status string, limit *int) ([]OrderInfo, error) {
	query, err := os.ReadFile(ordersPath + "id_by_status.sql")
	if err != nil {
		db.Log.Error("Database ShowByStatusOrders() -> Read SQL file", err)
		return nil, err
	}

	limitValue := sql.NullInt64{Valid: limit != nil}
	if limit != nil {
		limitValue.Int64 = int64(*limit)
	}

	rows, err := db.Query(string(query), status, limitValue)
	if err != nil {
		db.Log.Error("Database ShowByStatusOrders() -> tx.Query()", err)
		return nil, err
	}

	return db.readRowsOrderInfo(rows)
}

func (db *Database) RefundOrder(orderID int64) error {
	query, err := os.ReadFile(ordersPath + "refund_orders.sql")
	if err != nil {
		db.Log.Error("Database RefundOrder() -> Read SQL file", err)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		db.Log.Error("Database RefundOrder() -> db.Begin()", err)
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			db.Log.Error("Some troubles after tx.Rollback()", err)
		}
	}()

	if _, err = tx.Exec(string(query), orderID); err != nil {
		db.Log.Error("Database RefundOrder() -> tx.Exec()", err)
		return err
	}
	if err = tx.Commit(); err != nil {
		db.Log.Error("Database RefundOrder() -> tx.Commit()", err)
		return err
	}

	return nil
}

func (db *Database) readRowsOrder(rows *sql.Rows) ([]Order, error) {
	defer func() {
		if err := rows.Close(); err != nil {
			db.Log.Error("Some troubles after rows.Close()", err)
		}
	}()

	var products []Order

	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.OrderID, &o.CustomerID, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			db.Log.Error("Database readRowsOrder() -> parsing rows", err)
			return nil, err
		}
		products = append(products, o)
	}

	if err := rows.Err(); err != nil {
		db.Log.Error("Database readRowsOrder() -> rows.Err()", err)
		return nil, err
	}
	return products, nil
}

func (db *Database) ShowCustomerOrders(customerID int64, limit *int) ([]Order, error) {
	query, err := os.ReadFile(ordersPath + "show_customer_orders.sql")
	if err != nil {
		db.Log.Error("Database RefundOrder() -> Read SQL file", err)
		return nil, err
	}

	limitValue := sql.NullInt64{Valid: limit != nil}
	if limit != nil {
		limitValue.Int64 = int64(*limit)
	}

	rows, err := db.Query(string(query), customerID, limitValue)
	if err != nil {
		db.Log.Error("Database ShowByCustomerOrders() -> tx.Query()", err)
		return nil, err
	}

	return db.readRowsOrder(rows)
}

func (db *Database) readRowsOrderDetail(rows *sql.Rows) ([]OrderDetail, error) {
	defer func() {
		if err := rows.Close(); err != nil {
			db.Log.Error("Some troubles after rows.Close()", err)
		}
	}()

	var products []OrderDetail

	for rows.Next() {
		var o OrderDetail
		if err := rows.Scan(&o.OrderDetailID, &o.Quantity, &o.Price, &o.Name); err != nil {
			db.Log.Error("Database readRowsOrderDetail() -> parsing rows", err)
			return nil, err
		}
		products = append(products, o)
	}

	if err := rows.Err(); err != nil {
		db.Log.Error("Database readRowsOrderDetail() -> rows.Err()", err)
		return nil, err
	}
	return products, nil
}

func (db *Database) ShowOrderDetails(orderID int64, limit *int) ([]OrderDetail, error) {
	query, err := os.ReadFile(ordersPath + "show_order_details.sql")
	if err != nil {
		db.Log.Error("Database ShowOrderDetails() -> Read SQL file", err)
		return nil, err
	}

	limitValue := sql.NullInt64{Valid: limit != nil}
	if limit != nil {
		limitValue.Int64 = int64(*limit)
	}

	rows, err := db.Query(string(query), orderID, limitValue)
	if err != nil {
		db.Log.Error("Database ShowOrderDetails() -> tx.Query()", err)
		return nil, err
	}

	return db.readRowsOrderDetail(rows)
}

func (db *Database) UpdateStatusOrder(orderID int64, status string) error {
	query, err := os.ReadFile(ordersPath + "status_orders.sql")
	if err != nil {
		db.Log.Error("Database UpdateStatusOrder() -> Read SQL file", err)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		db.Log.Error("Database UpdateStatusOrder() -> db.Begin()", err)
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			db.Log.Error("Some troubles after tx.Rollback()", err)
		}
	}()

	if _, err = tx.Exec(string(query), orderID, status); err != nil {
		db.Log.Error("Database UpdateStatusOrder() -> tx.Exec()", err)
		return err
	}
	if err = tx.Commit(); err != nil {
		db.Log.Error("Database UpdateStatusOrder() -> tx.Commit()", err)
		return err
	}

	return nil
}

func (db *Database) ShowByStatusFullOrders(status string, limit *int) ([]Order, error) {
	query, err := os.ReadFile(ordersPath + "show_orders_by_status.sql")
	if err != nil {
		db.Log.Error("Database RefundOrder() -> Read SQL file", err)
		return nil, err
	}

	limitValue := sql.NullInt64{Valid: limit != nil}
	if limit != nil {
		limitValue.Int64 = int64(*limit)
	}

	rows, err := db.Query(string(query), status, limitValue)
	if err != nil {
		db.Log.Error("Database ShowByCustomerOrders() -> tx.Query()", err)
		return nil, err
	}

	return db.readRowsOrder(rows)
}

func (db *Database) CheckOrderExists(orderID int64) (bool, error) {
	query, err := os.ReadFile(ordersPath + "check_order_exists.sql")
	if err != nil {
		db.Log.Error("Database CheckOrderExists() -> Read SQL file", err)
		return false, err
	}

	var exists bool
	err = db.QueryRow(string(query), orderID).Scan(&exists)
	if err != nil {
		db.Log.Error("Database CheckOrderExists() -> QueryRow()", err)
		return false, err
	}

	return exists, nil
}

func (db *Database) GetOrderStatus(orderID int64) (string, error) {
	query, err := os.ReadFile(ordersPath + "check_order_exists.sql")
	if err != nil {
		db.Log.Error("Database GetOrderStatus() -> Read SQL file", err)
		return "", err
	}

	var status string
	err = db.QueryRow(string(query), orderID).Scan(&status)
	if err != nil {
		db.Log.Error("Database GetOrderStatus() -> Error retrieving order status", err)
		return "", err
	}
	return status, nil
}
