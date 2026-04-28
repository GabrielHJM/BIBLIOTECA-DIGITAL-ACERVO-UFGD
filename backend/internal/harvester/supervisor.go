package harvester

import (
	"biblioteca-digital-api/internal/pkg/logger"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

type APIStatus struct {
	Name      string
	IsOnline  bool
	LastCheck time.Time
	CheckURL  string
}

type Supervisor struct {
	mu     sync.RWMutex
	status map[string]*APIStatus
	client *http.Client
}

// Global instance of the supervisor
var GlobalSupervisor = NewSupervisor()

func NewSupervisor() *Supervisor {
	return &Supervisor{
		status: make(map[string]*APIStatus),
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// RegisterAPI adds a new API to be monitored
func (s *Supervisor) RegisterAPI(name, checkURL string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Assume online until proven otherwise
	s.status[name] = &APIStatus{
		Name:      name,
		IsOnline:  true,
		LastCheck: time.Now(),
		CheckURL:  checkURL,
	}
}

// IsOnline checks if a specific API is currently considered online
func (s *Supervisor) IsOnline(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if status, exists := s.status[name]; exists {
		return status.IsOnline
	}
	// If not registered, we assume it's online to not break unregistered APIs
	return true
}

// StartMonitoring begins the periodic check of all registered APIs
func (s *Supervisor) StartMonitoring(interval time.Duration) {
	logger.Info("Starting API Supervisor monitoring", zap.String("interval", interval.String()))
	go func() {
		for {
			s.checkAll()
			time.Sleep(interval)
		}
	}()
}

// checkAll iterates through all registered APIs and pings them
func (s *Supervisor) checkAll() {
	s.mu.RLock()
	// Create a copy to avoid holding the lock during HTTP requests
	apis := make(map[string]string)
	for name, status := range s.status {
		apis[name] = status.CheckURL
	}
	s.mu.RUnlock()

	for name, url := range apis {
		if url == "" {
			continue // Skip checking if no URL provided
		}

		online := s.ping(url)
		
		s.mu.Lock()
		if status, exists := s.status[name]; exists {
			if status.IsOnline != online {
				if online {
					logger.Info("API Restored", zap.String("api", name))
				} else {
					logger.Warn("API Down (Circuit Broken)", zap.String("api", name), zap.String("url", url))
				}
			}
			status.IsOnline = online
			status.LastCheck = time.Now()
		}
		s.mu.Unlock()
	}
}

// ping simply checks if the URL responds within the timeout and returns a 2xx or 3xx or 4xx status.
// We accept 4xx because it means the API is up, we might just be hitting it wrong without keys.
// We just want to avoid hanging on dead servers or 5xx errors.
func (s *Supervisor) ping(url string) bool {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		req, _ = http.NewRequest(http.MethodGet, url, nil)
	}
	
	// Fast timeout for health check
	resp, err := s.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return false
	}
	return true
}
