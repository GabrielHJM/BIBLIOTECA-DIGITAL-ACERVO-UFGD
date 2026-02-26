package harvester

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/logger"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// OAI-PMH Structures for parsing SciELO
type OAIRepository struct {
	XMLName     xml.Name `xml:"OAI-PMH"`
	ListRecords struct {
		Records []struct {
			Metadata struct {
				DublinCore struct {
					Titles       []string `xml:"title"`
					Creators     []string `xml:"creator"`
					Subjects     []string `xml:"subject"`
					Descriptions []string `xml:"description"`
					Publishers   []string `xml:"publisher"`
					Dates        []string `xml:"date"`
					Types        []string `xml:"type"`
					Identifiers  []string `xml:"identifier"`
					Languages    []string `xml:"language"`
				} `xml:"dc"`
			} `xml:"metadata"`
		} `xml:"record"`
		ResumptionToken string `xml:"resumptionToken"`
	} `xml:"ListRecords"`
}

type SciELOHarvester struct {
	BaseURL string
}

func NewSciELOHarvester() *SciELOHarvester {
	return &SciELOHarvester{
		BaseURL: "https://search.scielo.org/oai/scielo-oai.php",
	}
}

func (h *SciELOHarvester) Search(ctx context.Context, query string, category string, limit int) ([]material.Material, error) {
	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "science" // Default fallback
	}

	url := fmt.Sprintf("%s?verb=ListRecords&metadataPrefix=oai_dc", h.BaseURL)
	if searchTerm != "" {
		logger.Info("SciELO harvester: filtering by searchTerm", zap.String("query", searchTerm))
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("SciELO harvester: request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("scielo api error: %s", resp.Status)
	}

	var oai OAIRepository
	if err := xml.NewDecoder(resp.Body).Decode(&oai); err != nil {
		return nil, err
	}

	var materials []material.Material
	for i, rec := range oai.ListRecords.Records {
		if i >= limit && limit > 0 {
			break
		}

		dc := rec.Metadata.DublinCore
		if len(dc.Titles) == 0 {
			continue
		}

		if searchTerm != "" && !strings.Contains(strings.ToLower(dc.Titles[0]), strings.ToLower(searchTerm)) {
			// Fallback check: Se não achar no título, procura na listagem de publicadores ou descrição
			found := false
			for _, desc := range dc.Descriptions {
				if strings.Contains(strings.ToLower(desc), strings.ToLower(searchTerm)) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		year := 0
		if len(dc.Dates) > 0 {
			dateStr := dc.Dates[0]
			if len(dateStr) >= 4 {
				year, _ = strconv.Atoi(dateStr[:4])
			}
		}

		catName := category
		if catName == "" {
			catName = "Artigo Científico"
		}

		m := material.Material{
			Titulo:        dc.Titles[0],
			Autor:         strings.Join(dc.Creators, ", "),
			Descricao:     strings.Join(dc.Descriptions, " "),
			AnoPublicacao: year,
			Fonte:         "SciELO",
			Categoria:     catName,
			Disponivel:    true,
			CapaURL:       "https://images.unsplash.com/photo-1532012197267-da84d127e765?q=80&w=400",
		}

		for _, id := range dc.Identifiers {
			if strings.HasPrefix(id, "http") {
				if strings.Contains(id, ".pdf") {
					m.PDFURL = id
				} else if strings.Contains(id, "scielo.br") && !strings.Contains(id, "format=pdf") {
					// Se for um link do SciELO mas não for PDF, tentamos guardar como ExternoID
					// e futuramente podemos tentar inferir o PDF ou deixar o frontend tratar
					if m.ExternoID == "" {
						m.ExternoID = id
					}
					// Tentativa de conversão para link que o Google Viewer aceite melhor ou que seja PDF
					if strings.Contains(id, "sci_arttext") {
						m.PDFURL = id + "&format=pdf"
					}
				} else if m.ExternoID == "" {
					m.ExternoID = id
				}
			}
		}

		// Fallback: se não achou PDF mas tem ExternoID que parece ser landing page do SciELO
		if m.PDFURL == "" && strings.Contains(m.ExternoID, "scielo.br") {
			if strings.Contains(m.ExternoID, "sci_arttext") {
				m.PDFURL = m.ExternoID + "&format=pdf"
			}
		}


		materials = append(materials, m)
	}

	logger.Info("SciELO harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
