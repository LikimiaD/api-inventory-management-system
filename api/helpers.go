package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) respondWithError(w http.ResponseWriter, status int, error string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status": status,
		"error":  error,
	})
}

func (s *Server) respondWithToken(w http.ResponseWriter, token string) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status": http.StatusOK,
		"token":  token,
	})
}

func (s *Server) respondWithNew(w http.ResponseWriter, id int64) {
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status": http.StatusCreated,
		"id":     id,
	})
}

func (s *Server) respondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) respondNewOrder(w http.ResponseWriter, idOrder int64, idOrderDetail int64) {
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":          http.StatusCreated,
		"order_id":        idOrder,
		"order_detail_id": idOrderDetail,
	})
}

func (s *Server) respondWithStatus(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":  status,
		"message": message,
	})
}
