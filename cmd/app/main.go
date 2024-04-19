package main

import (
	"github.com/likimiad/golang-restapi-inventory-managment-system/api"
	"github.com/likimiad/golang-restapi-inventory-managment-system/iternal/config"
	"github.com/likimiad/golang-restapi-inventory-managment-system/iternal/database"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")

	db := database.InitDatabase(cfg.DatabaseConfig, log)
	db.InitTables()
	db.InitTrustedUsers()

	server := api.NewServer(log, db, cfg)
	if err := server.Start(cfg.HTTPServer.Address); err != nil {
		log.Error("Failed to start server", err)
		os.Exit(1)
	}
	log.Info("successfully start server")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
