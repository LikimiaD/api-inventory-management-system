package database

import (
	"database/sql"
	"errors"
	"os"
	"time"
)

type Product struct {
	ProductID   int64     `json:"product_id"`
	SupplierID  int64     `json:"supplier_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int64     `json:"quantity"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PurchaseRequest struct {
	ProductID    int64  `json:"product_id"`
	Name         string `json:"name"`
	SupplierID   int64  `json:"supplier_id"`
	ContactEmail string `json:"contact_email"`
}

var ErrNoProductFound = errors.New("no product found with the provided ID")

func (db *Database) AddProduct(supplierID int64, name, description string, price float64, quantity int64, category string) (int64, error) {
	query, err := os.ReadFile(productsPath + "add_products.sql")
	if err != nil {
		db.Log.Error("Database AddCustomer() -> Read SQL file", err)
		return 0, err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database AddCustomer() -> db.Begin()", err)
		return 0, err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	var productID int64
	err = tx.QueryRow(string(query), supplierID, name, description, price, quantity, category).Scan(&productID)
	if err != nil {
		db.Log.Error("Database AddCustomer() -> tx.QueryRow().Scan()", err)
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		db.Log.Error("Database AddCustomer() -> tx.Commit()", err)
		return 0, err
	}
	committed = true

	return productID, nil
}

func (db *Database) DeleteProduct(productID int64) error {
	query, err := os.ReadFile(productsPath + "delete_products.sql")
	if err != nil {
		db.Log.Error("Database DeleteProduct() -> Read SQL file", err)
		return err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database DeleteProduct() -> db.Begin()", err)
		return err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Error("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	result, err := tx.Exec(string(query), productID)
	if err != nil {
		db.Log.Error("Database DeleteProduct()", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.Log.Error("Database DeleteProduct() -> Checking rows affected", err)
		return err
	}
	if rowsAffected == 0 {
		return ErrNoProductFound
	}

	if err := tx.Commit(); err != nil {
		db.Log.Error("Database DeleteProduct() -> tx.Commit()", err)
		return err
	}
	committed = true

	return nil
}

func (db *Database) UpdateProduct(productID int64, name *string, supplierID *int64, description *string, price *float64, quantity *int64, category *string) error {
	query, err := os.ReadFile(productsPath + "set_products.sql")
	if err != nil {
		db.Log.Error("Database UpdateProduct() -> Read SQL file", err)
		return err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database UpdateProduct() -> db.Begin()", err)
		return err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	nameNull := sql.NullString{String: "", Valid: name != nil && *name != ""}
	supplierIDNull := sql.NullInt64{Int64: 0, Valid: supplierID != nil && *supplierID > 0}
	descriptionNull := sql.NullString{String: "", Valid: description != nil && *description != ""}
	priceNull := sql.NullFloat64{Float64: 0, Valid: price != nil && *price >= 0}
	quantityNull := sql.NullInt64{Int64: 0, Valid: quantity != nil && *quantity >= 0}
	categoryNull := sql.NullString{String: "", Valid: category != nil && *category != ""}

	if nameNull.Valid {
		nameNull.String = *name
	}
	if supplierIDNull.Valid {
		supplierIDNull.Int64 = *supplierID
	}
	if descriptionNull.Valid {
		descriptionNull.String = *description
	}
	if priceNull.Valid {
		priceNull.Float64 = *price
	}
	if quantityNull.Valid {
		quantityNull.Int64 = *quantity
	}
	if categoryNull.Valid {
		categoryNull.String = *category
	}

	result, err := tx.Exec(string(query), productID, nameNull, supplierIDNull, descriptionNull, priceNull, quantityNull, categoryNull)
	if err != nil {
		db.Log.Error("Database UpdateProduct() -> tx.Exec()", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.Log.Error("Database UpdateProduct() -> Checking rows affected", err)
		return err
	}
	if rowsAffected == 0 {
		return ErrNoProductFound
	}

	if err = tx.Commit(); err != nil {
		db.Log.Error("Database UpdateProduct() -> tx.Commit()", err)
		return err
	}
	committed = true

	return nil
}

func (db *Database) IDProduct(name string) (bool, error) {
	var exists bool

	query, err := os.ReadFile(productsPath + "id_products.sql")
	if err != nil {
		db.Log.Error("Database IDProduct() -> Read SQL file", err)
		return false, err
	}

	err = db.QueryRow(string(query), name).Scan(&exists)
	if err != nil {
		db.Log.Error("error checking if product exists:", err)
		return false, err
	}
	return exists, nil
}

func (db *Database) PurchaseRequestProducts(maxQuantity int64) ([]PurchaseRequest, error) {
	query, err := os.ReadFile(productsPath + "purchase_request_products.sql")
	if err != nil {
		db.Log.Error("Database PurchaseRequestProducts() -> Read SQL file", err)
		return nil, err
	}

	rows, err := db.Query(string(query), maxQuantity)
	if err != nil {
		db.Log.Error("Database PurchaseRequestProducts() -> tx.Query()", err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			db.Log.Error("Some troubles after rows.Close()", err)
		}
	}()

	var products []PurchaseRequest
	for rows.Next() {
		var p PurchaseRequest
		if err := rows.Scan(&p.ProductID, &p.Name, &p.SupplierID, &p.ContactEmail); err != nil {
			db.Log.Error("Database CheckProducts() -> parsing rows", err)
			return nil, err
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		db.Log.Error("Database CheckProducts() -> rows.Err()", err)
		return nil, err
	}
	return products, nil
}

func (db *Database) readRowsProduct(rows *sql.Rows) ([]Product, error) {
	defer func() {
		if err := rows.Close(); err != nil {
			db.Log.Error("Some troubles after rows.Close()", err)
		}
	}()

	var products []Product

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ProductID, &p.SupplierID, &p.Name, &p.Description, &p.Price, &p.Quantity, &p.Category, &p.CreatedAt, &p.UpdatedAt); err != nil {
			db.Log.Error("Database readRowsProduct() -> parsing rows", err)
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		db.Log.Error("Database readRowsProduct() -> rows.Err()", err)
		return nil, err
	}
	return products, nil
}

func (db *Database) ShowNotEmptyQuantityProducts(limit *int) ([]Product, error) {
	query, err := os.ReadFile(productsPath + "show_by_quantity_products.sql")
	if err != nil {
		db.Log.Error("Database ShowNotEmptyQuantityProducts() -> Read SQL file", err)
		return nil, err
	}

	limitValue := sql.NullInt64{Valid: limit != nil}
	if limit != nil {
		limitValue.Int64 = int64(*limit)
	}

	rows, err := db.Query(string(query), limitValue)
	if err != nil {
		db.Log.Error("Database ShowNotEmptyQuantityProducts() -> tx.Query()", err)
		return nil, err
	}

	return db.readRowsProduct(rows)
}

func (db *Database) ShowByCategoryProducts(category string, limit *int) ([]Product, error) {
	query, err := os.ReadFile(productsPath + "show_by_category_products.sql")
	if err != nil {
		db.Log.Error("Database ShowByCategoryProducts() -> Read SQL file", err)
		return nil, err
	}

	limitValue := sql.NullInt64{Valid: limit != nil}
	if limit != nil {
		limitValue.Int64 = int64(*limit)
	}

	rows, err := db.Query(string(query), category, limitValue)
	if err != nil {
		db.Log.Error("Database ShowByCategoryProducts() -> tx.Query()", err)
		return nil, err
	}

	return db.readRowsProduct(rows)
}

func (db *Database) ShowBetweenPriceProducts(min, max int64, limit *int) ([]Product, error) {
	query, err := os.ReadFile(productsPath + "show_price_products.sql")
	if err != nil {
		db.Log.Error("Database ShowBetweenPriceProducts() -> Read SQL file", err)
		return nil, err
	}

	limitValue := sql.NullInt64{Valid: limit != nil}
	if limit != nil {
		limitValue.Int64 = int64(*limit)
	}

	rows, err := db.Query(string(query), min, max, limitValue)
	if err != nil {
		db.Log.Error("Database ShowBetweenPriceProducts() -> tx.Query()", err)
		return nil, err
	}

	return db.readRowsProduct(rows)
}

func (db *Database) ShowProducts(limit *int) ([]Product, error) {
	query, err := os.ReadFile(productsPath + "show_products.sql")
	if err != nil {
		db.Log.Error("Database ShowProducts() -> Read SQL file", err)
		return nil, err
	}

	limitValue := sql.NullInt64{Valid: limit != nil}
	if limit != nil {
		limitValue.Int64 = int64(*limit)
	}

	rows, err := db.Query(string(query), limitValue)
	if err != nil {
		db.Log.Error("Database ShowProducts() -> db.Query()", err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			db.Log.Error("Some troubles after rows.Close()", err)
		}
	}()

	return db.readRowsProduct(rows)
}

func (db *Database) CheckProductAvailability(productID, quantity int64) (bool, error) {
	query, err := os.ReadFile(productsPath + "check_product_availability.sql")
	if err != nil {
		db.Log.Error("Database CheckProductAvailability() -> Read SQL file", err)
		return false, err
	}

	var currentQuantity int64
	err = db.QueryRow(string(query), productID).Scan(&currentQuantity)
	if err != nil {
		if err == sql.ErrNoRows {
			db.Log.Error("Database CheckProductAvailability() -> No product found with given ID", err)
			return false, nil
		}
		db.Log.Error("Database CheckProductAvailability() -> QueryRow()", err)
		return false, err
	}

	return currentQuantity >= quantity, nil
}

func (db *Database) CheckProductExists(productID int64) (bool, error) {
	query, err := os.ReadFile(productsPath + "check_product_exists.sql")
	if err != nil {
		db.Log.Error("Database CheckProductExists() -> Read SQL file", err)
		return false, err
	}

	var exists bool
	err = db.QueryRow(string(query), productID).Scan(&exists)
	if err != nil {
		db.Log.Error("Database CheckProductExists() -> Error checking if product exists", err)
		return false, err
	}

	return exists, nil
}
