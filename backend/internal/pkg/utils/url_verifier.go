package utils

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// URLVerifier handles efficient checking of URL availability
type URLVerifier struct {
	client *http.Client
}

// NewURLVerifier creates a new verifier with a custom client
func NewURLVerifier() *URLVerifier {
	return &URLVerifier{
		client: &http.Client{
			Timeout: 2 * time.Second,
			// Don't follow too many redirects and skip some body reading
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 3 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

// VerifyBatch checks multiple URLs in parallel
func (v *URLVerifier) VerifyBatch(ctx context.Context, urls []string) map[string]bool {
	results := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Semaphore to limit concurrency (e.g., max 10 parallel requests)
	sem := make(chan struct{}, 10)

	for _, url := range urls {
		if url == "" {
			continue
		}
		
		wg.Add(1)
		go func(targetURL string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			alive := v.IsAlive(ctx, targetURL)
			
			mu.Lock()
			results[targetURL] = alive
			mu.Unlock()
		}(url)
	}

	wg.Wait()
	return results
}

// IsAlive performs a rapid HEAD/GET request to check if a URL is reachable
func (v *URLVerifier) IsAlive(ctx context.Context, url string) bool {
	// 1. Try HEAD first (faster, no body)
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return true // Optimistic
	}
	
	// Add common headers to avoid being blocked as a bot
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := v.client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode < 400 {
			return true
		}
		if resp.StatusCode >= 500 {
			return true // Server error might be temporary
		}
	}

	// 2. Fallback to GET with a very small range if HEAD fails
	req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return true
	}
	req.Header.Set("Range", "bytes=0-0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err = v.client.Do(req)
	if err != nil {
		// On network/timeout/context errors, be optimistic
		return true
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		// Only definitively dead if 404/403 etc.
		return false
	}

	return true
}
