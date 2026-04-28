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

type PLOSResponse struct {
	Response struct {
		Docs []struct {
			ID              string   `json:"id"`
			Title           string   `json:"title_display"`
			Author          []string `json:"author_display"`
			Abstract        []string `json:"abstract"`
			PublicationDate string   `json:"publication_date"`
		} `json:"docs"`
	} `json:"response"`
}

type PLOSHarvester struct {
	BaseURL string
}

func NewPLOSHarvester() *PLOSHarvester {
	return &PLOSHarvester{
		BaseURL: "https://api.plos.org/search",
	}
}

func (h *PLOSHarvester) Search(ctx context.Context, query string, category string, limit int, offset int) ([]material.Material, error) {
	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "science"
	}

	searchURL := fmt.Sprintf("%s?q=title:%s&fl=id,title_display,author_display,abstract,publication_date&start=%d&rows=%d", h.BaseURL, url.QueryEscape(searchTerm), offset, limit)

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
		return nil, fmt.Errorf("plos api error: %s", resp.Status)
	}

	var data PLOSResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, doc := range data.Response.Docs {
		if doc.Title == "" {
			continue
		}

		autor := "Desconhecido"
		if len(doc.Author) > 0 {
			autor = strings.Join(doc.Author, ", ")
		}

		descricao := ""
		if len(doc.Abstract) > 0 {
			descricao = doc.Abstract[0]
		}

		year := 0
		if len(doc.PublicationDate) >= 4 {
			fmt.Sscanf(doc.PublicationDate[:4], "%d", &year)
		}

		pdfURL := fmt.Sprintf("https://journals.plos.org/plosone/article/file?id=%s&type=printable", doc.ID)

		m := material.Material{
			Titulo:        doc.Title,
			Autor:         autor,
			Descricao:     descricao,
			AnoPublicacao: year,
			Fonte:         "PLOS",
			Categoria:     category,
			ExternoID:     doc.ID,
			PDFURL:        pdfURL,
			Disponivel:    true,
			Dificuldade:   3,
			XP:            25,
			Relevancia:    10,
		}

		materials = append(materials, m)
	}

	logger.Info("PLOS harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
