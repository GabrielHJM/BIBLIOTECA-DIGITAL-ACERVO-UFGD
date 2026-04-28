package harvester

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/logger"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

type ArxivFeed struct {
	XMLName xml.Name     `xml:"feed"`
	Entries []ArxivEntry `xml:"entry"`
}

type ArxivEntry struct {
	ID        string `xml:"id"`
	Updated   string `xml:"updated"`
	Published string `xml:"published"`
	Title     string `xml:"title"`
	Summary   string `xml:"summary"`
	Authors   []struct {
		Name string `xml:"name"`
	} `xml:"author"`
	Links []struct {
		Href  string `xml:"href,attr"`
		Title string `xml:"title,attr"`
		Type  string `xml:"type,attr"`
	} `xml:"link"`
	Categories []struct {
		Term string `xml:"term,attr"`
	} `xml:"category"`
}

type ArXivHarvester struct {
	BaseURL string
}

func NewArXivHarvester() *ArXivHarvester {
	return &ArXivHarvester{
		BaseURL: "http://export.arxiv.org/api/query",
	}
}

func (h *ArXivHarvester) Search(ctx context.Context, query string, category string, limit int) ([]material.Material, error) {
	limiter := GetRateLimiter()

	searchTerm := query
	if searchTerm == "" {
		searchTerm = "all"
	}

	searchURL := fmt.Sprintf("%s?search_query=all:%s&start=0&max_results=%d", h.BaseURL, url.QueryEscape(searchTerm), limit)

	// Max 3 retries with exponential backoff
	var resp *http.Response
	var err error
	backoff := 1 * time.Second

	for i := 0; i < 3; i++ {
		// Wait for rate limiter (ArXiv: ~1 request every 3 seconds recommended, but we can try 1/s with retries)
		limiter.Wait(ctx, ProviderArXiv, 0.33, 1)

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
				logger.Warn("ArXiv rate limit hit, retrying...", zap.Int("attempt", i+1))
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			resp.Body.Close()
			return nil, fmt.Errorf("arxiv api error: %s", resp.Status)
		}

		if i < 2 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		if err != nil {
			logger.Error("ArXiv harvester: request failed", zap.Error(err))
			return nil, err
		}
		return nil, fmt.Errorf("arxiv api error after retries")
	}
	defer resp.Body.Close()

	var data ArxivFeed
	if err := xml.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, item := range data.Entries {
		if item.Title == "" {
			continue
		}

		var pdfURL string
		for _, link := range item.Links {
			if link.Title == "pdf" || link.Type == "application/pdf" {
				pdfURL = link.Href
				break
			} else if link.Type == "text/html" {
				pdfURL = link.Href
			}
		}

		if pdfURL == "" {
			if len(item.Links) > 0 {
				pdfURL = item.Links[0].Href
			} else {
				continue
			}
		}

		var authors []string
		for _, a := range item.Authors {
			authors = append(authors, a.Name)
		}

		year := 0
		if len(item.Published) >= 4 {
			fmt.Sscanf(item.Published[:4], "%d", &year)
		}

		// Gamification
		difficulty := 4 // ArXiv papers are usually dense
		xp := 10 + (difficulty * 5)
		relevance := 20

		cover := GetCoverFromGoogleBooks(item.Title, strings.Join(authors, ", "))

		m := material.Material{
			Titulo:        strings.ReplaceAll(item.Title, "\n", " "),
			Autor:         strings.Join(authors, ", "),
			Descricao:     strings.TrimSpace(item.Summary),
			AnoPublicacao: year,
			Fonte:         "ArXiv",
			Categoria:     category,
			ExternoID:     item.ID,
			CapaURL:       cover,
			PDFURL:        pdfURL,
			Disponivel:    true,
			Dificuldade:   difficulty,
			XP:            xp,
			Relevancia:    relevance,
		}

		if m.Categoria == "" {
			m.Categoria = "Pesquisa Avançada"
		}

		materials = append(materials, m)
	}

	logger.Info("ArXiv harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
