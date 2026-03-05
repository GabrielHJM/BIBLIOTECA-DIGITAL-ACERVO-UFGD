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

	"go.uber.org/zap"
)

type ArXivHarvester struct {
	BaseURL string
}

func NewArXivHarvester() *ArXivHarvester {
	return &ArXivHarvester{
		BaseURL: "http://export.arxiv.org/api/query",
	}
}

type ArXivEntry struct {
	ID        string `xml:"id"`
	Updated   string `xml:"updated"`
	Published string `xml:"published"`
	Title     string `xml:"title"`
	Summary   string `xml:"summary"`
	Author    []struct {
		Name string `xml:"name"`
	} `xml:"author"`
	Link []struct {
		Href string `xml:"href,attr"`
		Rel  string `xml:"rel,attr"`
		Type string `xml:"type,attr"`
	} `xml:"link"`
	Category []struct {
		Term string `xml:"term,attr"`
	} `xml:"category"`
}

type ArXivFeed struct {
	XMLName xml.Name     `xml:"feed"`
	Entry   []ArXivEntry `xml:"entry"`
}

func (h *ArXivHarvester) Search(ctx context.Context, query string, category string, limit int) ([]material.Material, error) {
	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "technology"
	}

	searchURL := fmt.Sprintf("%s?search_query=all:%s&start=0&max_results=%d", h.BaseURL, url.QueryEscape(searchTerm), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("ArXiv harvester: request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("arxiv api error: %s", resp.Status)
	}

	var feed ArXivFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("failed to decode arxiv xml: %w", err)
	}

	var materials []material.Material
	for _, entry := range feed.Entry {
		var authors []string
		for _, a := range entry.Author {
			authors = append(authors, a.Name)
		}

		year := 0
		if len(entry.Published) >= 4 {
			fmt.Sscanf(entry.Published[:4], "%d", &year)
		}

		// ArXiv links: usually one is for the abstract, another for the PDF
		var pdfURL string
		for _, link := range entry.Link {
			if link.Type == "application/pdf" || strings.Contains(link.Href, "pdf") {
				pdfURL = link.Href
				// Ensure PDF URL starts with https and is direct
				if !strings.HasPrefix(pdfURL, "http") {
					continue
				}
				// ArXiv links are usually http://arxiv.org/pdf/xxxx - upgrade to https if possible
				pdfURL = strings.Replace(pdfURL, "http://", "https://", 1)
				break
			}
		}

		// If no explicit PDF link found, try to derive it from the ID
		if pdfURL == "" && strings.Contains(entry.ID, "abs/") {
			pdfURL = strings.Replace(entry.ID, "abs/", "pdf/", 1) + ".pdf"
		} else if pdfURL == "" && strings.Contains(entry.ID, "arxiv.org/") {
			// Some IDs are just the URL
			pdfURL = entry.ID
			if !strings.HasSuffix(pdfURL, ".pdf") {
				pdfURL = strings.Replace(pdfURL, "abs/", "pdf/", 1) + ".pdf"
			}
		}

		if pdfURL == "" {
			continue
		}

		catName := category
		if catName == "" {
			if len(entry.Category) > 0 {
				catName = h.mapCategory(entry.Category[0].Term)
			} else {
				catName = "TECNOLOGIA"
			}
		}

		m := material.Material{
			Titulo:        strings.TrimSpace(entry.Title),
			Autor:         strings.Join(authors, ", "),
			Descricao:     strings.TrimSpace(entry.Summary),
			AnoPublicacao: year,
			Fonte:         "ArXiv",
			Categoria:     catName,
			ExternoID:     entry.ID,
			PDFURL:        pdfURL,
			Disponivel:    true,
			CapaURL:       GetCoverFromGoogleBooks(entry.Title, ""), // Use existing helper for cover
		}

		materials = append(materials, m)
	}

	logger.Info("ArXiv harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}

func (h *ArXivHarvester) mapCategory(term string) string {
	term = strings.ToLower(term)
	// Simple mapping for ArXiv categories to friendly portal categories
	if strings.HasPrefix(term, "cs.") {
		return "TECNOLOGIA"
	}
	if strings.HasPrefix(term, "math.") {
		return "MATEMÁTICA"
	}
	if strings.HasPrefix(term, "physics.") || strings.HasPrefix(term, "quant-ph") {
		return "CIÊNCIAS"
	}
	if strings.HasPrefix(term, "bio.") || strings.HasPrefix(term, "q-bio") {
		return "SAÚDE"
	}
	return "TECNOLOGIA" // Fallback for ArXiv which is mostly tech
}
