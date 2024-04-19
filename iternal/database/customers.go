package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

type Customer struct {
	ID      int64  `json:"customer_id"`
	Name    string `json:"customer_name"`
	Email   string `json:"customer_email"`
	Phone   string `json:"customer_phone"`
	Address string `json:"customer_address"`
}

var ErrNoCustomerFound = errors.New("no customer found with the provided ID")

func (db *Database) AddCustomer(name, email, phone, address string) (int64, error) {
	query, err := os.ReadFile(customersPath + "add_customers.sql")
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

	var customerID int64
	err = tx.QueryRow(string(query), name, email, phone, address).Scan(&customerID)
	if err != nil {
		db.Log.Error("Database AddCustomer() -> tx.QueryRow().Scan()", err)
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		db.Log.Error("Database AddCustomer() -> tx.Commit()", err)
		return 0, err
	}
	committed = true

	return customerID, nil
}

func (db *Database) CheckCustomer(email, phone string) (int, error) {
	query, err := os.ReadFile(customersPath + "check_by_email_customers.sql")
	if err != nil {
		db.Log.Error("Database CheckCustomer() -> Read SQL file", err)
		return 0, err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database CheckCustomer() -> db.Begin()", err)
		return 0, err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	var customerID int
	err = tx.QueryRow(string(query), email, phone).Scan(&customerID)
	if err != nil {
		db.Log.Error("Database CheckCustomer()", err)
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		db.Log.Error("Database CheckCustomer() -> tx.Commit()", err)
		return 0, err
	}
	committed = true

	return customerID, nil
}

func (db *Database) UpdateCustomer(customerID int64, name, email, phone, address *string) error {
	query, err := os.ReadFile(customersPath + "set_customers.sql")
	if err != nil {
		db.Log.Error("Database UpdateCustomer() -> Read SQL file", err)
		return err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database UpdateCustomer() -> db.Begin()", err)
		return err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	nameNull := sql.NullString{Valid: name != nil && *name != ""}
	if nameNull.Valid {
		nameNull.String = *name
	}
	emailNull := sql.NullString{Valid: email != nil && *email != ""}
	if emailNull.Valid {
		emailNull.String = *email
	}
	phoneNull := sql.NullString{Valid: phone != nil && *phone != ""}
	if phoneNull.Valid {
		phoneNull.String = *phone
	}
	addressNull := sql.NullString{Valid: address != nil && *address != ""}
	if addressNull.Valid {
		addressNull.String = *address
	}

	result, err := tx.Exec(string(query), customerID, nameNull, emailNull, phoneNull, addressNull)
	if err != nil {
		db.Log.Error("Database UpdateCustomer() -> tx.Exec()", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.Log.Error("Database UpdateCustomer() -> Checking rows affected", err)
		return err
	}
	if rowsAffected == 0 {
		return ErrNoCustomerFound
	}

	if err = tx.Commit(); err != nil {
		db.Log.Error("Database UpdateCustomer() -> tx.Commit()", err)
		return err
	}
	committed = true

	return nil
}

func (db *Database) DeleteCustomer(customerID int64) error {
	query, err := os.ReadFile(customersPath + "delete_customers.sql")
	if err != nil {
		db.Log.Error("Database DeleteCustomer() -> Read SQL file", err)
		return err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database DeleteCustomer() -> db.Begin()", err)
		return err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	result, err := tx.Exec(string(query), customerID)
	if err != nil {
		db.Log.Error("Database DeleteCustomer()", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.Log.Error("Database DeleteCustomer() -> Checking rows affected", err)
		return err
	}
	if rowsAffected == 0 {
		return ErrNoCustomerFound
	}

	if err := tx.Commit(); err != nil {
		db.Log.Error("Database DeleteCustomer() -> tx.Commit()", err)
		return err
	}
	committed = true

	return nil
}

func (db *Database) CheckEmailCustomer(contactEmail string) (int64, error) {
	query, err := os.ReadFile(customersPath + "check_by_email_customers.sql")
	if err != nil {
		db.Log.Error("Database CheckEmailCustomer() -> Read SQL file", err)
		return -1, err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database CheckEmailCustomer() -> db.Begin()", err)
		return -1, err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	var customerID int64
	err = tx.QueryRow(string(query), contactEmail).Scan(&customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		db.Log.Error("Database CheckEmailCustomer()", err)
		return -1, err
	}

	if err := tx.Commit(); err != nil {
		db.Log.Error("Database CheckEmailCustomer() -> tx.Commit()", err)
		return -1, err
	}
	committed = true

	return customerID, nil
}

func (db *Database) CheckPhoneCustomer(contactPhone string) (int64, error) {
	query, err := os.ReadFile(customersPath + "check_by_phone_customers.sql")
	if err != nil {
		db.Log.Error("Database CheckPhoneCustomer() -> Read SQL file", err)
		return -1, err
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		db.Log.Error("Database CheckPhoneCustomer() -> db.Begin()", err)
		return -1, err
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Warn("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	var customerID int64
	err = tx.QueryRow(string(query), contactPhone).Scan(&customerID)
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

	return customerID, nil
}

func (db *Database) readRowsCustomer(rows *sql.Rows) ([]Customer, error) {
	var customers []Customer
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address); err != nil {
			db.Log.Error("Database readRowsSupplier() -> rows.Scan()", err)
			return nil, err
		}
		customers = append(customers, c)
	}

	if err := rows.Err(); err != nil {
		db.Log.Error("Database readRowsSupplier() -> rows.Err()", err)
		return nil, err
	}
	return customers, nil
}

func (db *Database) ShowCustomers(limit *int) ([]Customer, error) {
	query, err := os.ReadFile(customersPath + "show_customers.sql")
	if err != nil {
		db.Log.Error("Database ShowCustomers() -> Read SQL file", err)
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
		db.Log.Error("Database ShowCustomers() -> db.Query()", err)
		return nil, err
	}
	defer rows.Close()

	return db.readRowsCustomer(rows)
}

func (db *Database) CheckCustomerExists(customerID int64) (bool, error) {
	query, err := os.ReadFile(customersPath + "check_customer_exists.sql")
	if err != nil {
		db.Log.Error("Database CheckCustomerExists() -> Read SQL file", err)
		return false, fmt.Errorf("error reading SQL file: %w", err)
	}

	var exists bool
	err = db.QueryRow(string(query), customerID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		db.Log.Error("Database CheckCustomerExists() -> Error checking if customer exists", err)
		return false, fmt.Errorf("error checking if customer exists: %w", err)
	}

	return exists, nil
}
