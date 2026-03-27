package harvester

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

type ProviderRateLimiter struct {
	limiters sync.Map
}

var globalRateLimiter = &ProviderRateLimiter{}

func GetRateLimiter() *ProviderRateLimiter {
	return globalRateLimiter
}

func (prl *ProviderRateLimiter) Wait(ctx context.Context, provider string, r rate.Limit, b int) error {
	limiter, _ := prl.limiters.LoadOrStore(provider, rate.NewLimiter(r, b))
	return limiter.(*rate.Limiter).Wait(ctx)
}

// Common providers
const (
	ProviderGoogleBooks    = "www.googleapis.com"
	ProviderSemanticScholar = "api.semanticscholar.org"
	ProviderArXiv          = "export.arxiv.org"
	ProviderCrossref       = "api.crossref.org"
	ProviderOpenLibrary    = "openlibrary.org"
	ProviderGutendex       = "gutendex.com"
)
