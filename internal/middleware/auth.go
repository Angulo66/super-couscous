package middleware

import (
	"go-basics/config"
	"net/http"

	"github.com/rs/zerolog/log"
)

const (
	BEARER_PREFIX = "Bearer "
)

func AuthMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != BEARER_PREFIX+cfg.Auth.Token {
			log.Error().Msg("Unauthorized request")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
