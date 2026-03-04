package middleware

import "net/http"

func Security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Proteção contra MIME-sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Proteção contra Clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Proteção contra XSS
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Referrer Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy (Básico)
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data: https:;")

		next.ServeHTTP(w, r)
	})
}
