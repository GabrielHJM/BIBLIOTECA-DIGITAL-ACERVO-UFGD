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

type GoogleBooksResponse struct {
	Items []struct {
		ID         string `json:"id"`
		VolumeInfo struct {
			Title         string   `json:"title"`
			Authors       []string `json:"authors"`
			PublishedDate string   `json:"publishedDate"`
			Description   string   `json:"description"`
			PageCount     int      `json:"pageCount"`
			Categories    []string `json:"categories"`
			ImageLinks    struct {
				Thumbnail string `json:"thumbnail"`
			} `json:"imageLinks"`
		} `json:"volumeInfo"`
		AccessInfo struct {
			Viewability string `json:"viewability"`
			Pdf         struct {
				IsAvailable  bool   `json:"isAvailable"`
				DownloadLink string `json:"downloadLink"`
				AcsTokenLink string `json:"acsTokenLink"`
			} `json:"pdf"`
			WebReaderLink string `json:"webReaderLink"`
		} `json:"accessInfo"`
	} `json:"items"`
}

type GoogleBooksHarvester struct {
	BaseURL string
}

func NewGoogleBooksHarvester() *GoogleBooksHarvester {
	return &GoogleBooksHarvester{
		BaseURL: "https://www.googleapis.com/books/v1/volumes",
	}
}

func (h *GoogleBooksHarvester) Search(ctx context.Context, query string, category string, limit int) ([]material.Material, error) {
	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "science"
	}

	apiKey := os.Getenv("GOOGLE_BOOKS_API_KEY")
	keyParam := ""
	if apiKey != "" {
		keyParam = "&key=" + apiKey
	}
	searchURL := fmt.Sprintf("%s?q=%s&filter=free-ebooks&langRestrict=pt&maxResults=%d%s", h.BaseURL, url.QueryEscape(searchTerm), limit, keyParam)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("GoogleBooks harvester: request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("googlebooks api error: %s", resp.Status)
	}

	var data GoogleBooksResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, item := range data.Items {
		if item.VolumeInfo.Title == "" {
			continue
		}

		pdfURL := ""
		// If Google explicitly says PDF is available and gives a download link
		if item.AccessInfo.Pdf.IsAvailable {
			if item.AccessInfo.Pdf.DownloadLink != "" {
				pdfURL = item.AccessInfo.Pdf.DownloadLink
			} else if item.AccessInfo.Pdf.AcsTokenLink != "" {
				pdfURL = item.AccessInfo.Pdf.AcsTokenLink
			}
		}

		if pdfURL == "" {
			continue
		}

		// Ensure the link leads to a PDF
		// We still allow googleapis.com links as they are the direct download source,
		// but we strip any non-pdf fallbacks like webReader
		lowerPDF := strings.ToLower(pdfURL)
		if !strings.HasSuffix(lowerPDF, ".pdf") && !strings.Contains(lowerPDF, "download") && !strings.Contains(lowerPDF, "acs") {
			continue
		}

		year := 0
		if len(item.VolumeInfo.PublishedDate) >= 4 {
			fmt.Sscanf(item.VolumeInfo.PublishedDate[:4], "%d", &year)
		}

		cover := item.VolumeInfo.ImageLinks.Thumbnail
		if cover != "" {
			cover = strings.ReplaceAll(cover, "http://", "https://")
		} else {
			cover = abstractAcademicCover(item.VolumeInfo.Title)
		}

		difficulty := 2
		if item.VolumeInfo.PageCount > 300 {
			difficulty = 3
		}
		if item.VolumeInfo.PageCount > 600 {
			difficulty = 4
		}

		xp := 10 + (difficulty * 5)
		relevance := 30 // High relevance for books

		cat := category
		if cat == "" {
			if len(item.VolumeInfo.Categories) > 0 {
				cat = item.VolumeInfo.Categories[0]
			} else {
				cat = "Livro"
			}
		}

		m := material.Material{
			Titulo:        item.VolumeInfo.Title,
			Autor:         strings.Join(item.VolumeInfo.Authors, ", "),
			Descricao:     item.VolumeInfo.Description,
			AnoPublicacao: year,
			Paginas:       item.VolumeInfo.PageCount,
			Fonte:         "Google Books",
			Categoria:     cat,
			ExternoID:     item.ID,
			CapaURL:       cover,
			PDFURL:        pdfURL,
			Disponivel:    true,
			Dificuldade:   difficulty,
			XP:            xp,
			Relevancia:    relevance,
		}

		materials = append(materials, m)
	}

	logger.Info("GoogleBooks harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
