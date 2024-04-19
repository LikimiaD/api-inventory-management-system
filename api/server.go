package api

import (
	"github.com/gorilla/mux"
	"github.com/likimiad/golang-restapi-inventory-managment-system/iternal/config"
	"github.com/likimiad/golang-restapi-inventory-managment-system/iternal/database"
	"log/slog"
	"net/http"
)

type Server struct {
	DB        *database.Database
	Router    *mux.Router
	Log       *slog.Logger
	SecretKey string
}

func NewServer(log *slog.Logger, db *database.Database, cfg *config.Config) *Server {
	server := &Server{
		DB:        db,
		Router:    mux.NewRouter(),
		Log:       log,
		SecretKey: cfg.SecretKey,
	}
	server.routes()
	return server
}

func (s *Server) Start(address string) error {
	s.Log.Info("Starting server", slog.String("address", address))
	return http.ListenAndServe(address, s.Router)
}
