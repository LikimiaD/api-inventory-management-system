package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type TrustedUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	var u TrustedUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Troubles with parsing data")
		return
	}

	if err := s.authenticateUser(w, u); err != nil {
		s.Log.Error("Authentication failed", slog.String("error", err.Error()))
	}
	s.Log.Info(fmt.Sprintf("User: %s created new token", u.Login))
}
