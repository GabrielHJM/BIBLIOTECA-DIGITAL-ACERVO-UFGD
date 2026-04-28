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

type InternetArchiveResponse struct {
	Response struct {
		Docs []struct {
			Identifier  string   `json:"identifier"`
			Title       string   `json:"title"`
			Creator     []string `json:"creator"`
			Description []string `json:"description"`
			Year        string   `json:"year"`
		} `json:"docs"`
	} `json:"response"`
}

type InternetArchiveHarvester struct {
	BaseURL string
}

func NewInternetArchiveHarvester() *InternetArchiveHarvester {
	return &InternetArchiveHarvester{
		BaseURL: "https://archive.org/advancedsearch.php",
	}
}

func (h *InternetArchiveHarvester) Search(ctx context.Context, query string, category string, limit int, offset int) ([]material.Material, error) {
	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "science" // default
	}

	// Calculate page (1-indexed for archive.org, though it accepts page parameter)
	page := (offset / limit) + 1

	q := fmt.Sprintf("(%s) AND mediatype:texts AND format:(Text PDF)", searchTerm)
	searchURL := fmt.Sprintf("%s?q=%s&fl[]=identifier,title,creator,description,year&output=json&rows=%d&page=%d", h.BaseURL, url.QueryEscape(q), limit, page)

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
		return nil, fmt.Errorf("internet archive api error: %s", resp.Status)
	}

	var data InternetArchiveResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, doc := range data.Response.Docs {
		if doc.Title == "" {
			continue
		}

		autor := "Unknown"
		if len(doc.Creator) > 0 {
			autor = strings.Join(doc.Creator, ", ")
		}

		descricao := ""
		if len(doc.Description) > 0 {
			descricao = doc.Description[0]
			if len(descricao) > 1000 {
				descricao = descricao[:1000] + "..."
			}
		}

		year := 0
		if doc.Year != "" {
			fmt.Sscanf(doc.Year, "%d", &year)
		}

		pdfURL := fmt.Sprintf("https://archive.org/download/%s/%s.pdf", doc.Identifier, doc.Identifier)
		coverURL := fmt.Sprintf("https://archive.org/services/img/%s", doc.Identifier)

		m := material.Material{
			Titulo:        doc.Title,
			Autor:         autor,
			Descricao:     descricao,
			AnoPublicacao: year,
			Fonte:         "Internet Archive",
			Categoria:     category,
			ExternoID:     doc.Identifier,
			CapaURL:       coverURL,
			PDFURL:        pdfURL,
			Disponivel:    true,
			Dificuldade:   2,
			XP:            20,
			Relevancia:    10,
		}

		if m.Categoria == "" {
			m.Categoria = "Arquivo Histórico"
		}

		materials = append(materials, m)
	}

	logger.Info("InternetArchive harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
