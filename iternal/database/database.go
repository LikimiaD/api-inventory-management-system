package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/likimiad/golang-restapi-inventory-managment-system/iternal/config"
	"log"
	"log/slog"
	"os"
	"strings"
)

const mainPath string = "sql/"
const customersPath string = mainPath + "customers/"
const suppliersPath = mainPath + "suppliers/"
const productsPath = mainPath + "products/"
const ordersPath = mainPath + "orders/"
const analyticsPath = mainPath + "analytics/"
const trustedUsersPath = mainPath + "trusted_users/"

type Database struct {
	*sql.DB
	Log *slog.Logger
}

func InitDatabase(cfg config.DatabaseConfig, slog *slog.Logger) *Database {
	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name)
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		log.Fatalf("error connection to database: %s", err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("error during connection check: %s", err.Error())
	}

	slog.Info("successfully connect to the database")
	return &Database{db, slog}
}

func (db *Database) InitTables() {
	file, err := os.ReadFile(mainPath + "create_tables.sql")

	if err != nil {
		log.Fatalf("error reading the database loading script: %s", err.Error())
	}

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		_, err := db.Exec(request)
		if err != nil {
			log.Fatalf("error database during table initialization: %s", err.Error())
		}
	}
	db.Log.Info("successful check/initialization of tables")
}
