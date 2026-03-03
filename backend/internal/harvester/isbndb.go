package harvester

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"go.uber.org/zap"
)

type ISBNdbResponse struct {
	Books []struct {
		Title         string   `json:"title"`
		Authors       []string `json:"authors"`
		DatePublished string   `json:"date_published"`
		Isbn13        string   `json:"isbn13"`
		Image         string   `json:"image"`
		Synopsis      string   `json:"synopsis"`
		Subjects      []string `json:"subjects"`
	} `json:"books"`
}

type ISBNdbHarvester struct {
	BaseURL string
	APIKey  string
}

func NewISBNdbHarvester() *ISBNdbHarvester {
	apiKey := os.Getenv("ISBNDB_API_KEY") // Placeholder para a chave da API
	return &ISBNdbHarvester{
		BaseURL: "https://api2.isbndb.com/books",
		APIKey:  apiKey,
	}
}

func (h *ISBNdbHarvester) Search(ctx context.Context, query string, category string, limit int) ([]material.Material, error) {
	if h.APIKey == "" {
		// Log and return empty gracefully if no API key is set
		logger.Warn("ISBNdb harvester: API key not set, skipping. Defina a variável de ambiente ISBNDB_API_KEY.")
		return []material.Material{}, nil
	}

	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "science"
	}

	searchURL := fmt.Sprintf("%s/%s?page=1&pageSize=%d", h.BaseURL, url.PathEscape(searchTerm), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", h.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("ISBNdb harvester: request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("isbndb api error: %s", resp.Status)
	}

	var data ISBNdbResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, item := range data.Books {
		if item.Title == "" {
			continue
		}

		catName := category
		if catName == "" {
			catName = "Livro"
		}

		year := 0
		if len(item.DatePublished) >= 4 {
			fmt.Sscanf(item.DatePublished[:4], "%d", &year)
		}

		tags := []string{}
		if len(item.Subjects) > 0 {
			for i, s := range item.Subjects {
				if i < 5 {
					tags = append(tags, s)
				}
			}
		}

		m := material.Material{
			Titulo:        item.Title,
			Autor:         strings.Join(item.Authors, ", "),
			Descricao:     item.Synopsis,
			AnoPublicacao: year,
			Fonte:         "ISBNdb",
			Categoria:     catName,
			ExternoID:     item.Isbn13,
			ISBN:          item.Isbn13,
			CapaURL:       item.Image,
			Disponivel:    true,
			Tags:          tags,
		}

		materials = append(materials, m)
	}

	logger.Info("ISBNdb harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
