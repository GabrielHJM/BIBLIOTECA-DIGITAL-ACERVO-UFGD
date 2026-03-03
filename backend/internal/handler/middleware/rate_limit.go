package middleware

import (
	"net/http"
	"sync"
	"time"

	"biblioteca-digital-api/internal/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// client representa um visitante pela perspectiva do limitador de taxa
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

func init() {
	go cleanupClients()
}

func cleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

// RateLimit é um middleware que limita o número de requisições por IP
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		// Simplificação: usando RemoteAddr (em produção, usar X-Forwarded-For se atrás de um proxy/Cloudflare)
		if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			ip = forwardedFor
		}

		mu.Lock()
		if _, found := clients[ip]; !found {
			// Permite 5 requests por segundo com um burst de 10
			clients[ip] = &client{limiter: rate.NewLimiter(5, 10)}
		}
		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			logger.Warn("Rate limit exceeded", zap.String("ip", ip))
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
