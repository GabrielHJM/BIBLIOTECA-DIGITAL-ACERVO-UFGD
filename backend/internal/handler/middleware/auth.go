package middleware

import (
	"biblioteca-digital-api/pkg/auth"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UsuarioIDKey contextKey = "usuario_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token não fornecido", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
			return
		}

		userID, err := auth.VerifyToken(parts[1])
		if err != nil || userID == 0 {
			http.Error(w, "Token inválido ou expirado", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UsuarioIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
