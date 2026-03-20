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

func (h *SemanticScholarHarvester) Search(ctx context.Context, query string, category string, limit int) ([]material.Material, error) {
	time.Sleep(1500 * time.Millisecond) // Prevents HTTP 429 Too Many Requests

	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "science"
	}

	searchURL := fmt.Sprintf("%s?query=%s&limit=%d&fields=title,authors,year,abstract,openAccessPdf,citationCount,journal", h.BaseURL, url.QueryEscape(searchTerm), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("SemanticScholar harvester: request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("semanticscholar api error: %s", resp.Status)
	}

	var data SemanticScholarResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, item := range data.Data {
		if item.Title == "" || item.OpenAccessPdf.URL == "" {
			continue
		}

		// STRICT PDF RULE
		pdfURL := item.OpenAccessPdf.URL
		if !strings.HasSuffix(strings.ToLower(strings.Split(pdfURL, "?")[0]), ".pdf") {
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
