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
	mu       sync.Mutex
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients sync.Map
)

func init() {
	go cleanupClients()
}

func cleanupClients() {
	for {
		time.Sleep(time.Minute)
		clients.Range(func(key, value interface{}) bool {
			ip := key.(string)
			c := value.(*client)

			c.mu.Lock()
			if time.Since(c.lastSeen) > 3*time.Minute {
				clients.Delete(ip)
			}
			c.mu.Unlock()
			return true
		})
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

		value, _ := clients.LoadOrStore(ip, &client{
			limiter:  rate.NewLimiter(5, 10),
			lastSeen: time.Now(),
		})
		c := value.(*client)

		c.mu.Lock()
		c.lastSeen = time.Now()
		allowed := c.limiter.Allow()
		c.mu.Unlock()

		if !allowed {
			logger.Warn("Rate limit exceeded", zap.String("ip", ip))
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

