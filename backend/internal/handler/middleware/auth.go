package middleware

import "net/http"

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // lógica de autenticação
        next.ServeHTTP(w, r)
    })
}
