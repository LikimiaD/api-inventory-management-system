package database

import (
	"log"
	"os"
	"strings"
)

func (db *Database) InitTrustedUsers() {
	db.Log.Info("Check for trusted users")
	if db.IsEmptyTrustedUsers() {
		db.Log.Info("Register trusted users in the database")
		db.LoadTrustedUsers()
	} else {
		db.Log.Info("Trusted users already exist in the database")
	}
}

func (db *Database) IsEmptyTrustedUsers() bool {
	var countUsers int64
	err := db.QueryRow(`SELECT COUNT(*) FROM trusted_users;`).Scan(&countUsers)
	if err != nil {
		log.Fatalf("Troubles with base check of count trusted users: %s", err)
	}
	return countUsers == 0
}

func (db *Database) LoadTrustedUsers() {
	file, err := os.ReadFile(trustedUsersPath + "base_add_trusted_users.sql")

	if err != nil {
		log.Fatalf("error reading the trusted users loading script: %s", err)
	}

	if string(file) == "" {
		log.Fatalf("error, \"base_add_trusted_users.sql\" can't be empty")
	}

	tx, err := db.Begin()
	committed := false
	if err != nil {
		log.Fatalf("LoadTrustedUsers() -> db.Begin: %s", err)
	}
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				db.Log.Error("Some troubles after tx.Rollback()", err)
			}
		}
	}()

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		_, err := tx.Exec(request)
		if err != nil {
			log.Fatalf("error database during table initialization: %s", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("LoadTrustedUsers() -> db.Commit: %s", err)
	}
	committed = true
}

func (db *Database) CheckTrustedUser(login string) (bool, string, error) {
	query, err := os.ReadFile(trustedUsersPath + "check_trusted_user.sql")
	if err != nil {
		db.Log.Error("Database CheckTrustedUser() -> Read SQL file", err)
		return false, "", err
	}
	var password string
	if err := db.QueryRow(string(query), login).Scan(&password); err != nil {
		db.Log.Error("Database CheckTrustedUser() -> db.QueryRow()", err)
	}
	return password != "", password, nil
}
