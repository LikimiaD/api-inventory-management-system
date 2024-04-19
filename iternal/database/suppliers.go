package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

var ErrNoSupplierFound = errors.New("no supplier found with the provided ID")

type Supplier struct {
	SupplierID   int64  `json:"id"`
	Name         string `json:"name"`
	ContactName  string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
}

func (db *Database) AddSupplier(name, contactName, contactEmail, contactPhone string) (int64, error) {
	query, err := os.ReadFile(suppliersPath + "add_suppliers.sql")
	if err != nil {
		db.Log.Error("Database AddSupplier() -> Read SQL file", err)
		return 0, err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database AddSupplier() -> db.Begin()", err)
		return -1, err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	var supplierID int64
	err = tx.QueryRow(string(query), name, contactName, contactEmail, contactPhone).Scan(&supplierID)
	if err != nil {
		db.Log.Error("Database AddSupplier()", err)
		return -1, err
	}

	if err := tx.Commit(); err != nil {
		db.Log.Error("Database AddSupplier() -> tx.Commit()", err)
		return -1, err
	}
	committed = true

	return supplierID, nil
}

func (db *Database) CheckEmailSupplier(contactEmail string) (int64, error) {
	query, err := os.ReadFile(suppliersPath + "check_by_email_suppliers.sql")
	if err != nil {
		db.Log.Error("Database CheckEmailSupplier() -> Read SQL file", err)
		return -1, err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database CheckEmailSupplier() -> db.Begin()", err)
		return -1, err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	var supplierID int64
	err = tx.QueryRow(string(query), contactEmail).Scan(&supplierID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		db.Log.Error("Database CheckEmailSupplier()", err)
		return -1, err
	}

	if err := tx.Commit(); err != nil {
		db.Log.Error("Database CheckEmailSupplier() -> tx.Commit()", err)
		return -1, err
	}
	committed = true

	return supplierID, nil
}

func (db *Database) CheckPhoneSupplier(contactPhone string) (int64, error) {
	query, err := os.ReadFile(suppliersPath + "check_by_phone_suppliers.sql")
	if err != nil {
		db.Log.Error("Database CheckPhoneSupplier() -> Read SQL file", err)
		return -1, err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database CheckPhoneSupplier() -> db.Begin()", err)
		return -1, err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	var supplierID int64
	err = tx.QueryRow(string(query), contactPhone).Scan(&supplierID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		db.Log.Error("Database CheckPhoneSupplier()", err)
		return -1, err
	}

	if err := tx.Commit(); err != nil {
		db.Log.Error("Database CheckPhoneSupplier() -> tx.Commit()", err)
		return -1, err
	}
	committed = true

	return supplierID, nil
}

func (db *Database) DeleteSupplier(supplierID int64) error {
	query, err := os.ReadFile(suppliersPath + "delete_suppliers.sql")
	if err != nil {
		db.Log.Error("Database DeleteSupplier() -> Read SQL file", err)
		return err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database DeleteSupplier() -> db.Begin()", err)
		return err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	result, err := tx.Exec(string(query), supplierID)
	if err != nil {
		db.Log.Error("Database DeleteSupplier()", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.Log.Error("Database DeleteSupplier() -> Checking rows affected", err)
		return err
	}

	if rowsAffected == 0 {
		db.Log.Info(fmt.Sprintf("Database DeleteSupplier() -> No supplier found with ID %d", supplierID))
		return ErrNoSupplierFound
	}

	if err := tx.Commit(); err != nil {
		db.Log.Error("Database DeleteSupplier() -> tx.Commit()", err)
		return err
	}
	committed = true

	return nil
}

func (db *Database) UpdateSupplier(supplierID int64, name, contactName, contactEmail, contactPhone *string) error {
	query, err := os.ReadFile(suppliersPath + "set_suppliers.sql")
	if err != nil {
		db.Log.Error("Database UpdateSupplier() -> Read SQL file", err)
		return err
	}
	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database UpdateSupplier() -> db.Begin()", err)
		return err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	nameNull := sql.NullString{String: "", Valid: false}
	contactNameNull := sql.NullString{String: "", Valid: false}
	contactEmailNull := sql.NullString{String: "", Valid: false}
	contactPhoneNull := sql.NullString{String: "", Valid: false}

	if name != nil {
		nameNull = sql.NullString{String: *name, Valid: true}
	}
	if contactName != nil {
		contactNameNull = sql.NullString{String: *contactName, Valid: true}
	}
	if contactEmail != nil {
		contactEmailNull = sql.NullString{String: *contactEmail, Valid: true}
	}
	if contactPhone != nil {
		contactPhoneNull = sql.NullString{String: *contactPhone, Valid: true}
	}

	result, err := tx.Exec(string(query), supplierID, nameNull, contactNameNull, contactEmailNull, contactPhoneNull)
	if err != nil {
		db.Log.Error("Database UpdateSupplier() -> tx.Exec()", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.Log.Error("Database UpdateSupplier() -> Checking rows affected", err)
		return err
	}

	if rowsAffected == 0 {
		return ErrNoSupplierFound
	}

	if err = tx.Commit(); err != nil {
		db.Log.Error("Database UpdateSupplier() -> tx.Commit()", err)
		return err
	}
	committed = true

	return nil
}

func (db *Database) readRowsSupplier(rows *sql.Rows) ([]Supplier, error) {
	var suppliers []Supplier
	for rows.Next() {
		var s Supplier
		if err := rows.Scan(&s.SupplierID, &s.Name, &s.ContactName, &s.ContactEmail, &s.ContactPhone); err != nil {
			db.Log.Error("Database readRowsSupplier() -> rows.Scan()", err)
			return nil, err
		}
		suppliers = append(suppliers, s)
	}

	if err := rows.Err(); err != nil {
		db.Log.Error("Database readRowsSupplier() -> rows.Err()", err)
		return nil, err
	}
	return suppliers, nil
}

func (db *Database) ShowSuppliers(limit *int) ([]Supplier, error) {
	query, err := os.ReadFile(suppliersPath + "show_suppliers.sql")
	if err != nil {
		db.Log.Error("Database ShowSuppliers() -> Read SQL file", err)
		return nil, err
	}

	var limitValue sql.NullInt64
	if limit != nil {
		limitValue = sql.NullInt64{Int64: int64(*limit), Valid: true}
	} else {
		limitValue = sql.NullInt64{Int64: 1000, Valid: true}
	}

	rows, err := db.Query(string(query), limitValue)
	if err != nil {
		db.Log.Error("Database ShowSuppliers() -> db.Query()", err)
		return nil, err
	}
	defer rows.Close()

	return db.readRowsSupplier(rows)
}

func (db *Database) CheckSupplierExists(supplierID int64) (bool, error) {
	query, err := os.ReadFile(productsPath + "check_supplier_exists.sql")
	if err != nil {
		db.Log.Error("Database CheckSupplierExists() -> Read SQL file", err)
		return false, err
	}

	var exists bool
	err = db.QueryRow(string(query), supplierID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		db.Log.Error("Database CheckSupplierExists() -> Error checking if supplier exists", err)
		return false, err
	}

	return exists, nil
}
