package harvester

import (
	"biblioteca-digital-api/internal/domain/material"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// OpenAlex Harvester
type OpenAlexHarvester struct { client *http.Client }
func NewOpenAlexHarvester() *OpenAlexHarvester { return &OpenAlexHarvester{client: &http.Client{Timeout: 10 * time.Second}} }
func (h *OpenAlexHarvester) Search(ctx context.Context, query, category string, limit, offset int) ([]material.Material, error) {
	apiURL := fmt.Sprintf("https://api.openalex.org/works?search=%s&per-page=%d&page=%d", url.QueryEscape(query), limit, (offset/limit)+1)
	return h.fetch(ctx, apiURL, "OpenAlex")
}

// Zenodo Harvester
type ZenodoHarvester struct { client *http.Client }
func NewZenodoHarvester() *ZenodoHarvester { return &ZenodoHarvester{client: &http.Client{Timeout: 10 * time.Second}} }
func (h *ZenodoHarvester) Search(ctx context.Context, query, category string, limit, offset int) ([]material.Material, error) {
	apiURL := fmt.Sprintf("https://zenodo.org/api/records?q=%s&size=%d&page=%d", url.QueryEscape(query), limit, (offset/limit)+1)
	return h.fetch(ctx, apiURL, "Zenodo")
}

// HAL Harvester
type HALHarvester struct { client *http.Client }
func NewHALHarvester() *HALHarvester { return &HALHarvester{client: &http.Client{Timeout: 10 * time.Second}} }
func (h *HALHarvester) Search(ctx context.Context, query, category string, limit, offset int) ([]material.Material, error) {
	apiURL := fmt.Sprintf("https://api.archives-ouvertes.fr/search/?q=%s&rows=%d&start=%d", url.QueryEscape(query), limit, offset)
	return h.fetch(ctx, apiURL, "HAL")
}

// PubMed Harvester
type PubMedHarvester struct { client *http.Client }
func NewPubMedHarvester() *PubMedHarvester { return &PubMedHarvester{client: &http.Client{Timeout: 10 * time.Second}} }
func (h *PubMedHarvester) Search(ctx context.Context, query, category string, limit, offset int) ([]material.Material, error) {
	apiURL := fmt.Sprintf("https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pmc&term=%s&retmode=json&retmax=%d&retstart=%d", url.QueryEscape(query), limit, offset)
	return h.fetch(ctx, apiURL, "PubMed")
}

// OSF Harvester
type OSFHarvester struct { client *http.Client }
func NewOSFHarvester() *OSFHarvester { return &OSFHarvester{client: &http.Client{Timeout: 10 * time.Second}} }
func (h *OSFHarvester) Search(ctx context.Context, query, category string, limit, offset int) ([]material.Material, error) {
	apiURL := fmt.Sprintf("https://api.osf.io/v2/nodes/?filter[title]=%s&page[size]=%d", url.QueryEscape(query), limit)
	return h.fetch(ctx, apiURL, "OSF")
}

// BASE Harvester (Stub)
type BASEHarvester struct { client *http.Client }
func NewBASEHarvester() *BASEHarvester { return &BASEHarvester{client: &http.Client{Timeout: 10 * time.Second}} }
func (h *BASEHarvester) Search(ctx context.Context, query, category string, limit, offset int) ([]material.Material, error) {
	return []material.Material{}, nil // Requires IP auth
}

// CORE Harvester (Stub)
type COREHarvester struct { client *http.Client }
func NewCOREHarvester() *COREHarvester { return &COREHarvester{client: &http.Client{Timeout: 10 * time.Second}} }
func (h *COREHarvester) Search(ctx context.Context, query, category string, limit, offset int) ([]material.Material, error) {
	return []material.Material{}, nil // Requires API Key
}

// SciELO Harvester (Stub)
type SciELOHarvester struct { client *http.Client }
func NewSciELOHarvester() *SciELOHarvester { return &SciELOHarvester{client: &http.Client{Timeout: 10 * time.Second}} }
func (h *SciELOHarvester) Search(ctx context.Context, query, category string, limit, offset int) ([]material.Material, error) {
	return []material.Material{}, nil // No public JSON endpoint easily searchable
}

// Generic Fetch function for the new harvesters
func (h *OpenAlexHarvester) fetch(ctx context.Context, apiURL, source string) ([]material.Material, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	resp, err := h.client.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { return nil, fmt.Errorf("bad status: %d", resp.StatusCode) }

	var mats []material.Material
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil { return nil, err }

	// Tenta extrair base genérica
	var results []interface{}
	if r, ok := data["results"].([]interface{}); ok { results = r } else if h, ok := data["hits"].(map[string]interface{}); ok {
		if hr, ok := h["hits"].([]interface{}); ok { results = hr }
	} else if r, ok := data["data"].([]interface{}); ok { results = r }

	for _, item := range results {
		iMap, ok := item.(map[string]interface{})
		if !ok { continue }
		
		title := "Sem Título"
		if t, ok := iMap["title"].(string); ok { title = t } else if t, ok := iMap["display_name"].(string); ok { title = t }
		
		if title == "Sem Título" || title == "" { continue }

		link := ""
		if o, ok := iMap["open_access"].(map[string]interface{}); ok {
			if u, ok := o["oa_url"].(string); ok { link = u }
		}

		mats = append(mats, material.Material{
			Titulo: title,
			Descricao: fmt.Sprintf("Material acadêmico extraído de %s", source),
			Categoria: "Ciência",
			PDFURL: link,
			Fonte: source,
			Disponivel: true,
			Status: "aprovado",
		})
	}
	return mats, nil
}

func (h *ZenodoHarvester) fetch(ctx context.Context, apiURL, source string) ([]material.Material, error) {
	return (&OpenAlexHarvester{client: h.client}).fetch(ctx, apiURL, source)
}
func (h *HALHarvester) fetch(ctx context.Context, apiURL, source string) ([]material.Material, error) {
	return (&OpenAlexHarvester{client: h.client}).fetch(ctx, apiURL, source)
}
func (h *PubMedHarvester) fetch(ctx context.Context, apiURL, source string) ([]material.Material, error) {
	return (&OpenAlexHarvester{client: h.client}).fetch(ctx, apiURL, source)
}
func (h *OSFHarvester) fetch(ctx context.Context, apiURL, source string) ([]material.Material, error) {
	return (&OpenAlexHarvester{client: h.client}).fetch(ctx, apiURL, source)
}
