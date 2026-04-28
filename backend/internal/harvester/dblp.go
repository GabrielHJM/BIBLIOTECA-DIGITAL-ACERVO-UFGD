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

	"go.uber.org/zap"
)

type DBLPResponse struct {
	Result struct {
		Hits struct {
			Hit []struct {
				Info struct {
					Authors struct {
						Author []interface{} `json:"author"`
					} `json:"authors"`
					Title string `json:"title"`
					Year  string `json:"year"`
					Ee    string `json:"ee"`
					Key   string `json:"key"`
				} `json:"info"`
			} `json:"hit"`
		} `json:"hits"`
	} `json:"result"`
}

type DBLPHarvester struct {
	BaseURL string
}

func NewDBLPHarvester() *DBLPHarvester {
	return &DBLPHarvester{
		BaseURL: "https://dblp.org/search/publ/api",
	}
}

func (h *DBLPHarvester) Search(ctx context.Context, query string, category string, limit int, offset int) ([]material.Material, error) {
	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "computer science"
	}

	searchURL := fmt.Sprintf("%s?q=%s&format=json&h=%d&f=%d", h.BaseURL, url.QueryEscape(searchTerm), limit, offset)

	req, errReq := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if errReq != nil {
		return nil, errReq
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dblp api error: %s", resp.Status)
	}

	var data DBLPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, hit := range data.Result.Hits.Hit {
		info := hit.Info
		if info.Title == "" || info.Ee == "" {
			continue
		}

		// Parse authors (DBLP can return author as object or array of objects/strings depending on count)
		autor := "Unknown"
		if info.Authors.Author != nil {
			autor = "Vários Autores"
		}

		year := 0
		if info.Year != "" {
			fmt.Sscanf(info.Year, "%d", &year)
		}

		// The "ee" link is usually the DOI or publisher link, but we'll capture it as the source URL.
		// It might not always be a direct PDF, but it's the official open link usually.
		pdfURL := info.Ee
		if !strings.HasPrefix(pdfURL, "http") {
			continue
		}

		m := material.Material{
			Titulo:        info.Title,
			Autor:         autor,
			Descricao:     "Publicação da base DBLP de Ciência da Computação.",
			AnoPublicacao: year,
			Fonte:         "DBLP",
			Categoria:     category,
			ExternoID:     info.Key,
			PDFURL:        pdfURL,
			Disponivel:    true,
			Dificuldade:   4,
			XP:            25,
			Relevancia:    10,
		}

		materials = append(materials, m)
	}

	logger.Info("DBLP harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
