package middleware

import (
	"biblioteca-digital-api/internal/pkg/logger"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Logger é um middleware que registra detalhes de cada requisição HTTP
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Serve a requisição
		next.ServeHTTP(w, r)

		// Loga os detalhes após a execução em formato estruturado
		logger.Info("Incoming Request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Duration("duration", time.Since(start)),
		)
	})
}
