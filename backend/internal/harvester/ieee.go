package harvester

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// IEEEHarvester implementa a coleta de dados da API do IEEE Xplore.
// Você poderá precisar de uma chave de API para o IEEE se o acesso não for aberto para determinados campos, mas
// utilizaremos a URL pública/open access conforme possível, ou simular o comportamento padrão.
// Referência da documentação do IEEE Xplore API (Open Access/Metadata API).
type IEEEHarvester struct {
	BaseURL string
	APIKey  string // Placeholder caso o usuário precise injetar uma chave de API futuramente
	Client  *http.Client
}

func NewIEEEHarvester() *IEEEHarvester {
	return &IEEEHarvester{
		BaseURL: "http://ieeexploreapi.ieee.org/api/v1/search/articles", // API pública base, necessita de apikey em prod
		Client:  &http.Client{},
	}
}

func (h *IEEEHarvester) Search(ctx context.Context, query string, category string, limit int) ([]material.Material, error) {
	if query == "" {
		if category != "" {
			query = category
		} else {
			query = "technology" // Busca padrão se tudo estiver vazio
		}
	}

	// Como IEEE é uma base forte em tecnologia e exatas, podemos buscar os dados.
	// NOTA: A API do IEEE exige API Key para chamadas reais.
	// Se a API Key não for fornecida e houver erro de Auth, o sistema deve logar e continuar silenciosamente.
	// Vamos construir a request.

	reqURL, err := url.Parse(h.BaseURL)
	if err != nil {
		return nil, err
	}

	q := reqURL.Query()
	q.Add("querytext", query)
	q.Add("max_records", strconv.Itoa(limit))
	// q.Add("apikey", h.APIKey) // Necessário em ambiente real

	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		logger.Error("Erro ao consultar IEEE Xplore", zap.Error(err))
		return nil, err // Retorna nil, err para não falhar toda a interface do multi\_harvester
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Warn("API do IEEE Xplore retornou status não OK", zap.Int("status", resp.StatusCode))
		// Simulando comportamento gracioso sem API Key
		return []material.Material{}, nil
	}

	var apiResp struct {
		Articles []struct {
			Title   string `json:"title"`
			Authors struct {
				Authors []struct {
					FullName string `json:"full_name"`
				} `json:"authors"`
			} `json:"authors"`
			Abstract        string `json:"abstract"`
			PublicationYear string `json:"publication_year"`
			PdfURL          string `json:"pdf_url"`
			DOI             string `json:"doi"`
			ArticleNumber   string `json:"article_number"`
		} `json:"articles"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON do IEEE: %w", err)
	}

	var materiais []material.Material
	for _, item := range apiResp.Articles {
		autoresStr := ""
		for i, a := range item.Authors.Authors {
			if i > 0 {
				autoresStr += "; "
			}
			autoresStr += a.FullName
		}
		if autoresStr == "" {
			autoresStr = "IEEE Xplore"
		}

		ano := 0
		if item.PublicationYear != "" {
			ano, _ = strconv.Atoi(item.PublicationYear)
		}

		// Constrói PDF URL se disponível e aberto, senão joga pra página do artigo
		pdfUrl := item.PdfURL
		if pdfUrl == "" && item.ArticleNumber != "" {
			pdfUrl = fmt.Sprintf("https://ieeexplore.ieee.org/document/%s", item.ArticleNumber)
		}

		catName := category
		if catName == "" {
			catName = "Tecnologia"
		}

		materiais = append(materiais, material.Material{
			Titulo:        strings.TrimSpace(item.Title),
			Autor:         autoresStr,
			Descricao:     strings.TrimSpace(item.Abstract),
			AnoPublicacao: ano,
			Categoria:     catName,
			CapaURL:       "", // IEEE geralmente não fornece capa de livro, é artigo
			PDFURL:        pdfUrl,
			ExternoID:     item.DOI,
			Fonte:         "IEEE",
		})
	}

	logger.Info("Resultados da busca no IEEE processados", zap.Int("quantidade", len(materiais)))
	return materiais, nil
}
