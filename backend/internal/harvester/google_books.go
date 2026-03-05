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

type GoogleBooksHarvester struct {
	BaseURL string
}

func NewGoogleBooksHarvester() *GoogleBooksHarvester {
	return &GoogleBooksHarvester{
		BaseURL: "https://www.googleapis.com/books/v1/volumes",
	}
}

type GoogleBooksResponse struct {
	Items []struct {
		ID         string `json:"id"`
		VolumeInfo struct {
			Title               string   `json:"title"`
			Authors             []string `json:"authors"`
			Description         string   `json:"description"`
			PublishedDate       string   `json:"publishedDate"`
			PageCount           int      `json:"pageCount"`
			Categories          []string `json:"categories"`
			IndustryIdentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentifiers"`
			ImageLinks struct {
				Thumbnail  string `json:"thumbnail"`
				Large      string `json:"large"`
				ExtraLarge string `json:"extraLarge"`
			} `json:"imageLinks"`
		} `json:"volumeInfo"`
		AccessInfo struct {
			WebReaderLink string `json:"webReaderLink"`
			Pdf           struct {
				IsAvailable  bool   `json:"isAvailable"`
				DownloadLink string `json:"downloadLink"`
				AcsTokenLink string `json:"acsTokenLink"`
			} `json:"pdf"`
		} `json:"accessInfo"`
	} `json:"items"`
}

func GetCoverFromGoogleBooks(title string, author string) string {
	query := title
	if author != "" {
		query += "+inauthor:" + author
	}
	searchURL := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=%s&maxResults=1", url.QueryEscape(query))

	resp, err := http.Get(searchURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var data GoogleBooksResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ""
	}

	if len(data.Items) > 0 {
		img := data.Items[0].VolumeInfo.ImageLinks
		// Try to get the highest quality possible, fallback to thumbnail. Google Books thumbnails are HTTP by default, upgrade to HTTPS
		cover := img.ExtraLarge
		if cover == "" {
			cover = img.Large
		}
		if cover == "" {
			cover = img.Thumbnail
		}
		return strings.ReplaceAll(cover, "http://", "https://")
	}

	return ""
}

func (h *GoogleBooksHarvester) Search(ctx context.Context, query string, category string, limit int) ([]material.Material, error) {
	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "science"
	}

	// Filter for free-ebooks to maximize chance of having a PDF link available
	searchURL := fmt.Sprintf("%s?q=%s&filter=free-ebooks&maxResults=%d", h.BaseURL, url.QueryEscape(searchTerm), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("Google Books harvester: request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google books api error: %s", resp.Status)
	}

	var data GoogleBooksResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, item := range data.Items {
		vol := item.VolumeInfo
		acc := item.AccessInfo

		if vol.Title == "" {
			continue
		}

		// Ensure strictly there's a PDF available.
		// Google Books sometimes provides acsTokenLink for DRM-protected PDFs, we really want direct PDFs if possible,
		// but `isAvailable` usually means there's some PDF form.
		if !acc.Pdf.IsAvailable {
			continue
		}

		pdfURL := acc.Pdf.DownloadLink
		if pdfURL == "" {
			pdfURL = acc.Pdf.AcsTokenLink
		}
		if pdfURL == "" {
			pdfURL = acc.WebReaderLink
		}

		// Ensure we don't return files that clearly aren't PDFs if a download link is provided
		if pdfURL != "" && !strings.HasSuffix(strings.ToLower(strings.Split(pdfURL, "?")[0]), ".pdf") {
			// Let's still trust the IsAvailable flag if there is no explicit extension in the URL as Google APIs often use dynamic endpoints
			if !strings.Contains(strings.ToLower(pdfURL), "pdf") && !strings.Contains(strings.ToLower(pdfURL), "books.google.com") {
				continue
			}
		}

		if pdfURL == "" {
			continue
		}

		year := 0
		if len(vol.PublishedDate) >= 4 {
			fmt.Sscanf(vol.PublishedDate[:4], "%d", &year)
		}

		isbn := ""
		for _, id := range vol.IndustryIdentifiers {
			if id.Type == "ISBN_13" || id.Type == "ISBN_10" {
				isbn = id.Identifier
				break
			}
		}

		cover := vol.ImageLinks.ExtraLarge
		if cover == "" {
			cover = vol.ImageLinks.Large
		}
		if cover == "" {
			cover = vol.ImageLinks.Thumbnail
		}
		cover = strings.ReplaceAll(cover, "http://", "https://")

		catName := category
		if catName == "" {
			if len(vol.Categories) > 0 {
				catName = vol.Categories[0]
			} else {
				catName = "Livro"
			}
		}

		m := material.Material{
			Titulo:        vol.Title,
			Autor:         strings.Join(vol.Authors, ", "),
			Descricao:     vol.Description,
			AnoPublicacao: year,
			ISBN:          isbn,
			Fonte:         "Google Books",
			Categoria:     catName,
			ExternoID:     item.ID,
			Paginas:       vol.PageCount,
			CapaURL:       cover,
			PDFURL:        pdfURL,
			Disponivel:    true,
		}

		materials = append(materials, m)
	}

	logger.Info("Google Books harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
