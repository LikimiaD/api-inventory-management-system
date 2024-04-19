package api

import (
	"fmt"
	"net/http"
)

func (s *Server) authenticateUser(w http.ResponseWriter, u TrustedUser) error {
	if u.Login == "" || u.Password == "" {
		s.respondWithError(w, http.StatusBadRequest, "Login or password is empty")
		return fmt.Errorf("empty credentials")
	}

	userExists, correctPassword, err := s.DB.CheckTrustedUser(u.Login)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return err
	}

	if !userExists {
		s.respondWithError(w, http.StatusBadRequest, "User does not exist")
		return fmt.Errorf("non-existent user")
	}

	if correctPassword != u.Password {
		s.respondWithError(w, http.StatusBadRequest, "Incorrect password")
		return fmt.Errorf("wrong password")
	}

	token, err := s.generateJWT(u)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return err
	}

	s.respondWithToken(w, token)
	return nil
}
