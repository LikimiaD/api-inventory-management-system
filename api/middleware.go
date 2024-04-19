package api

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

func (s *Server) isAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			s.respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(s.SecretKey), nil
		})

		if err != nil || !token.Valid {
			s.respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) generateJWT(u TrustedUser) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = u.Login
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(s.SecretKey))

	if err != nil {
		s.Log.Error("error routes side -> generatedJWT()", err)
		return "", err
	}

	return tokenString, nil
}

func (s *Server) getUserFromToken(r *http.Request) (string, error) {
	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
	if tokenString == "" {
		return "", fmt.Errorf("no Authorization header provided")
	}

	// Проверяем префикс "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return "", fmt.Errorf("authorization header must start with Bearer")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.SecretKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user, ok := claims["user"].(string)
		if !ok {
			return "", fmt.Errorf("error retrieving user from token")
		}
		return user, nil
	}

	return "", fmt.Errorf("invalid token")
}
