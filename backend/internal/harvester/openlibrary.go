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

type OpenLibraryResponse struct {
	NumFound int `json:"numFound"`
	Docs     []struct {
		Key             string   `json:"key"`
		Title           string   `json:"title"`
		AuthorName      []string `json:"author_name"`
		FirstPublishYear int      `json:"first_publish_year"`
		Subject         []string `json:"subject"`
		CoverI          int      `json:"cover_i"`
		HasFulltext     bool     `json:"has_fulltext"`
		IA              []string `json:"ia"` // Internet Archive IDs
	} `json:"docs"`
}

type OpenLibraryHarvester struct {
	BaseURL string
}

func NewOpenLibraryHarvester() *OpenLibraryHarvester {
	return &OpenLibraryHarvester{
		BaseURL: "https://openlibrary.org/search.json",
	}
}

func (h *OpenLibraryHarvester) Search(ctx context.Context, query string, category string, page int, limit int) ([]material.Material, error) {
	limiter := GetRateLimiter()

	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "science"
	}

	searchURL := fmt.Sprintf("%s?q=%s+language:por&has_fulltext=true&page=%d&limit=%d", h.BaseURL, url.QueryEscape(searchTerm), page, limit)

	// Max 3 retries
	var resp *http.Response
	var err error
	backoff := 1 * time.Second

	for i := 0; i < 3; i++ {
		// Wait for rate limiter (Open Library: Fair use policy, ~1 request per second)
		limiter.Wait(ctx, ProviderOpenLibrary, 1, 1)

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
				logger.Warn("OpenLibrary rate limit hit, retrying...", zap.Int("attempt", i+1))
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			resp.Body.Close()
			return nil, fmt.Errorf("openlibrary api error: %s", resp.Status)
		}

		if i < 2 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		if err != nil {
			logger.Error("OpenLibrary harvester: request failed", zap.Error(err))
			return nil, err
		}
		return nil, fmt.Errorf("openlibrary api error after retries")
	}
	defer resp.Body.Close()

	var data OpenLibraryResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, doc := range data.Docs {
		if doc.Title == "" {
			continue
		}

		// We ONLY want materials that have an Internet Archive ID for direct reading
		if len(doc.IA) == 0 {
			continue
		}

		// Internet Archive PDF direct link (most browsers render this online)
		pdfURL := fmt.Sprintf("https://archive.org/download/%s/%s.pdf", doc.IA[0], doc.IA[0])
		
		// If we wanted the reader: pdfURL = fmt.Sprintf("https://archive.org/details/%s/mode/2up", doc.IA[0])
		// But the user specifically mentioned "pdf do navegador" (browser PDF)
		if pdfURL == "" {
			continue
		}

		cover := ""
		if doc.CoverI > 0 {
			cover = fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-M.jpg", doc.CoverI)
		} else {
			cover = abstractAcademicCover(doc.Title)
		}

		cat := category
		if cat == "" && len(doc.Subject) > 0 {
			cat = doc.Subject[0]
		}
		if cat == "" {
			cat = "Livro"
		}

		materials = append(materials, material.Material{
			Titulo:        doc.Title,
			Autor:         strings.Join(doc.AuthorName, ", "),
			AnoPublicacao: doc.FirstPublishYear,
			Fonte:         "Open Library",
			Categoria:     cat,
			ExternoID:     doc.Key,
			CapaURL:       cover,
			PDFURL:        pdfURL,
			Disponivel:    true,
			Dificuldade:   3,
			XP:            25,
			Relevancia:    15,
		})
	}

	return materials, nil
}
