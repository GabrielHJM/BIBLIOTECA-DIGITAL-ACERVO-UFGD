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

type GutendexResponse struct {
	Count   int `json:"count"`
	Results []struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Authors []struct {
			Name      string `json:"name"`
			BirthYear int    `json:"birth_year"`
			DeathYear int    `json:"death_year"`
		} `json:"authors"`
		Subjects  []string           `json:"subjects"`
		Formats   map[string]string `json:"formats"`
		DownloadCount int            `json:"download_count"`
	} `json:"results"`
}

type GutendexHarvester struct {
	BaseURL string
}

func NewGutendexHarvester() *GutendexHarvester {
	return &GutendexHarvester{
		BaseURL: "https://gutendex.com/books",
	}
}

func (h *GutendexHarvester) Search(ctx context.Context, query string, category string, page int) ([]material.Material, error) {
	limiter := GetRateLimiter()

	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "classic"
	}

	searchURL := fmt.Sprintf("%s?search=%s&languages=pt&page=%d", h.BaseURL, url.QueryEscape(searchTerm), page)

	// Max 3 retries
	var resp *http.Response
	var err error
	backoff := 1 * time.Second

	for i := 0; i < 3; i++ {
		// Gutendex is fairly open, 2 requests per second is safe
		limiter.Wait(ctx, ProviderGutendex, 2, 1)

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
				logger.Warn("Gutendex rate limit hit, retrying...", zap.Int("attempt", i+1))
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			resp.Body.Close()
			return nil, fmt.Errorf("gutendex api error: %s", resp.Status)
		}

		if i < 2 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		if err != nil {
			logger.Error("Gutendex harvester: request failed", zap.Error(err))
			return nil, err
		}
		return nil, fmt.Errorf("gutendex api error after retries")
	}
	defer resp.Body.Close()

	var data GutendexResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, res := range data.Results {
		if res.Title == "" {
			continue
		}

		// Find the best reading link (Prioritizing HTML for online reading over downloads)
		pdfURL := ""
		if url, ok := res.Formats["text/html; charset=utf-8"]; ok {
			pdfURL = url
		} else if url, ok := res.Formats["application/pdf"]; ok {
			pdfURL = url
		} else if url, ok := res.Formats["text/plain; charset=utf-8"]; ok {
			pdfURL = url
		} else if url, ok := res.Formats["application/epub+zip"]; ok {
			pdfURL = url
		}

		if pdfURL == "" {
			continue
		}

		var authors []string
		for _, a := range res.Authors {
			authors = append(authors, a.Name)
		}

		cover := ""
		if url, ok := res.Formats["image/jpeg"]; ok {
			cover = url
		} else {
			cover = abstractAcademicCover(res.Title)
		}

		cat := category
		if cat == "" && len(res.Subjects) > 0 {
			cat = res.Subjects[0]
		}
		if cat == "" {
			cat = "Clássico"
		}

		materials = append(materials, material.Material{
			Titulo:        res.Title,
			Autor:         strings.Join(authors, ", "),
			Fonte:         "Project Gutenberg",
			Categoria:     cat,
			ExternoID:     fmt.Sprintf("gut:%d", res.ID),
			CapaURL:       cover,
			PDFURL:        pdfURL,
			Disponivel:    true,
			Dificuldade:   4, // Classics can be hard
			XP:            30,
			Relevancia:    res.DownloadCount / 100,
		})
	}

	return materials, nil
}
