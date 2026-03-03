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

type OpenLibraryResponse struct {
	Docs []struct {
		Title            string   `json:"title"`
		AuthorName       []string `json:"author_name"`
		FirstPublishYear int      `json:"first_publish_year"`
		Key              string   `json:"key"`
		CoverI           int      `json:"cover_i"`
		Subject          []string `json:"subject"`
		Isbn             []string `json:"isbn"`
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

func (h *OpenLibraryHarvester) Search(ctx context.Context, query string, category string, limit int) ([]material.Material, error) {
	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "books"
	}

	searchURL := fmt.Sprintf("%s?q=%s&limit=%d", h.BaseURL, url.QueryEscape(searchTerm), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("Open Library harvester: request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("open library api error: %s", resp.Status)
	}

	var data OpenLibraryResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, item := range data.Docs {
		if item.Title == "" {
			continue
		}

		catName := category
		if catName == "" {
			catName = "Livro"
		}

		coverURL := ""
		if item.CoverI > 0 {
			coverURL = fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-L.jpg", item.CoverI)
		}

		isbn := ""
		if len(item.Isbn) > 0 {
			isbn = item.Isbn[0]
		}

		tags := []string{}
		if len(item.Subject) > 0 {
			for i, s := range item.Subject {
				if i < 5 { // limit to 5 tags
					tags = append(tags, s)
				}
			}
		}

		m := material.Material{
			Titulo:        item.Title,
			Autor:         strings.Join(item.AuthorName, ", "),
			AnoPublicacao: item.FirstPublishYear,
			Fonte:         "Open Library",
			Categoria:     catName,
			ExternoID:     item.Key,                             // e.g., /works/OL12345W
			PDFURL:        "https://openlibrary.org" + item.Key, // Link to the book page
			CapaURL:       coverURL,
			ISBN:          isbn,
			Disponivel:    true,
			Tags:          tags,
		}

		materials = append(materials, m)
	}

	logger.Info("Open Library harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
