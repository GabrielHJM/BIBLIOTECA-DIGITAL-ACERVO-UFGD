package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger é um middleware que registra detalhes de cada requisição HTTP
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Serve a requisição
		next.ServeHTTP(w, r)

		// Loga os detalhes após a execução
		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}
