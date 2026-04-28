package harvester

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

type SemanticScholarResponse struct {
	Data []struct {
		PaperID       string `json:"paperId"`
		Title         string `json:"title"`
		Abstract      string `json:"abstract"`
		Year          int    `json:"year"`
		CitationCount int    `json:"citationCount"`
		Authors       []struct {
			Name string `json:"name"`
		} `json:"authors"`
		OpenAccessPdf struct {
			URL string `json:"url"`
		} `json:"openAccessPdf"`
		Journal struct {
			Name string `json:"name"`
		} `json:"journal"`
	} `json:"data"`
}

type SemanticScholarHarvester struct {
	BaseURL string
}

func NewSemanticScholarHarvester() *SemanticScholarHarvester {
	return &SemanticScholarHarvester{
		BaseURL: "https://api.semanticscholar.org/graph/v1/paper/search",
	}
}

func (h *SemanticScholarHarvester) Search(ctx context.Context, query string, category string, limit int, offset int) ([]material.Material, error) {
	limiter := GetRateLimiter()

	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "science"
	}

	searchURL := fmt.Sprintf("%s?query=%s&limit=%d&offset=%d&fields=title,authors,year,abstract,openAccessPdf,citationCount,journal", h.BaseURL, url.QueryEscape(searchTerm), limit, offset)

	// Max 3 retries with exponential backoff
	var resp *http.Response
	var err error
	backoff := 1 * time.Second

	for i := 0; i < 3; i++ {
		// Wait for rate limiter (Semantic Scholar: ~1 request per second for free tier)
		limiter.Wait(ctx, ProviderSemanticScholar, 1, 1)

		req, errReq := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
		if errReq != nil {
			return nil, errReq
		}

		resp, err = http.DefaultClient.Do(req)
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				break
			}
			if resp.StatusCode == http.StatusTooManyRequests {
				resp.Body.Close()
				logger.Warn("SemanticScholar rate limit hit, retrying...", zap.Int("attempt", i+1))
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			resp.Body.Close()
			return nil, fmt.Errorf("semanticscholar api error: %s", resp.Status)
		}

		if i < 2 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		if err != nil {
			logger.Error("SemanticScholar harvester: request failed", zap.Error(err))
			return nil, err
		}
		return nil, fmt.Errorf("semanticscholar api error after retries")
	}
	defer resp.Body.Close()

	var data SemanticScholarResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, item := range data.Data {
		if item.Title == "" || item.OpenAccessPdf.URL == "" {
			continue
		}

		// Relaxed PDF RULE
		pdfURL := item.OpenAccessPdf.URL
		lowerURL := strings.ToLower(pdfURL)
		if !strings.HasSuffix(strings.Split(lowerURL, "?")[0], ".pdf") && 
		   !strings.Contains(lowerURL, "pdf") && 
		   !strings.Contains(lowerURL, "openaccess") {
			continue
		}

		var authors []string
		for _, a := range item.Authors {
			authors = append(authors, a.Name)
		}

		// Gamification based on CitationCount
		difficulty := 1
		if item.CitationCount > 100 {
			difficulty = 5
		} else if item.CitationCount > 50 {
			difficulty = 4
		} else if item.CitationCount > 20 {
			difficulty = 3
		} else if item.CitationCount > 5 {
			difficulty = 2
		}

		xp := 10 + (difficulty * 5)

		// Fetch cover - Try Title first, then Journal fallback
		cover := GetCoverFromGoogleBooks(item.Title, strings.Join(authors, ", "))
		if cover == "" && item.Journal.Name != "" {
			cover = GetCoverFromGoogleBooks(item.Journal.Name, "")
		}

		m := material.Material{
			Titulo:        item.Title,
			Autor:         strings.Join(authors, ", "),
			Descricao:     item.Abstract,
			AnoPublicacao: item.Year,
			Fonte:         "Semantic Scholar",
			Categoria:     category,
			ExternoID:     item.PaperID,
			CapaURL:       cover,
			PDFURL:        pdfURL,
			Disponivel:    true,
			Dificuldade:   difficulty,
			XP:            xp,
			Relevancia:    item.CitationCount,
		}

		if m.Categoria == "" {
			m.Categoria = "Artigo Científico"
		}

		materials = append(materials, m)
	}

	logger.Info("SemanticScholar harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
