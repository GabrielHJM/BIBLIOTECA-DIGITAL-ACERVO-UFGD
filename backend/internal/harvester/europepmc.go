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

type EuropePMCResponse struct {
	ResultList struct {
		Result []struct {
			ID           string `json:"id"`
			Pmid         string `json:"pmid"`
			Pmcid        string `json:"pmcid"`
			Title        string `json:"title"`
			AuthorString string `json:"authorString"`
			JournalTitle string `json:"journalTitle"`
			PubYear      string `json:"pubYear"`
			AbstractText string `json:"abstractText"`
		} `json:"result"`
	} `json:"resultList"`
}

type EuropePMCHarvester struct {
	BaseURL string
}

func NewEuropePMCHarvester() *EuropePMCHarvester {
	return &EuropePMCHarvester{
		BaseURL: "https://www.ebi.ac.uk/europepmc/webservices/rest/search",
	}
}

func (h *EuropePMCHarvester) Search(ctx context.Context, query string, category string, limit int, offset int) ([]material.Material, error) {
	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "medicine"
	}

	// Europe PMC pagination uses Cursor, but for basic offset we can use cursorMark=* for first page 
	// Or since we just want random pages, we can append random sort or use a basic search.
	// We'll use cursorMark=* and just pull the top results for the specific query.
	// To fake offset, we can append year to the query if needed, but let's stick to the basic query.
	q := fmt.Sprintf("%s AND HAS_PDF:y", searchTerm)
	
	searchURL := fmt.Sprintf("%s?query=%s&format=json&resultType=core&pageSize=%d", h.BaseURL, url.QueryEscape(q), limit)

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
		return nil, fmt.Errorf("europe pmc api error: %s", resp.Status)
	}

	var data EuropePMCResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, item := range data.ResultList.Result {
		if item.Title == "" || item.Pmcid == "" {
			continue
		}

		year := 0
		if item.PubYear != "" {
			fmt.Sscanf(item.PubYear, "%d", &year)
		}

		pdfURL := fmt.Sprintf("https://europepmc.org/articles/%s?pdf=render", item.Pmcid)

		m := material.Material{
			Titulo:        item.Title,
			Autor:         item.AuthorString,
			Descricao:     strings.TrimSpace(item.AbstractText),
			AnoPublicacao: year,
			Fonte:         "Europe PMC",
			Categoria:     category,
			ExternoID:     item.Pmcid,
			PDFURL:        pdfURL,
			Disponivel:    true,
			Dificuldade:   4,
			XP:            30,
			Relevancia:    15,
		}

		materials = append(materials, m)
	}

	logger.Info("EuropePMC harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
