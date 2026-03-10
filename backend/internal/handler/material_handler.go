package handler

import (
	"biblioteca-digital-api/internal/harvester"
	"biblioteca-digital-api/internal/pkg/cache"
	"biblioteca-digital-api/internal/repository"
	"biblioteca-digital-api/internal/usecase/material"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func RegisterMaterialRoutes(mux *http.ServeMux, db *sql.DB, c cache.Cache) {
	repo := &repository.MaterialPostgres{DB: db}
	mh := harvester.NewMultiSourceHarvester()

	listarUC := &material.ListarConteudosUseCase{Repo: repo, Harvester: mh, Cache: c}
	buscarUC := &material.BuscarMaterialUseCase{Repo: repo}
	similaresUC := &material.BuscarSimilaresUseCase{Repo: repo}
	pesquisarUC := &material.PesquisarMaterialUseCase{Repo: repo, Harvester: mh, Cache: c}
	favoritarUC := &material.FavoritarMaterialUseCase{Repo: repo}
	historicoUC := &material.HistoricoLeituraUseCase{Repo: repo}
	avaliarUC := &material.AvaliarMaterialUseCase{Repo: repo}

	mux.HandleFunc("GET /materiais", func(w http.ResponseWriter, r *http.Request) {
		termo := r.URL.Query().Get("q")
		categoria := r.URL.Query().Get("categoria")
		fonte := r.URL.Query().Get("fonte")
		anoInicioStr := r.URL.Query().Get("ano_inicio")
		var anoInicio int
		if anoInicioStr != "" {
			var err error
			anoInicio, err = strconv.Atoi(anoInicioStr)
			if err != nil || anoInicio < 0 {
				JSONError(w, "Ano de início inválido", http.StatusBadRequest)
				return
			}
		}

		anoFimStr := r.URL.Query().Get("ano_fim")
		var anoFim int
		if anoFimStr != "" {
			var err error
			anoFim, err = strconv.Atoi(anoFimStr)
			if err != nil || anoFim < 0 {
				JSONError(w, "Ano de fim inválido", http.StatusBadRequest)
				return
			}
		}

		limitStr := r.URL.Query().Get("limit")
		limit := 10
		if limitStr != "" {
			var err error
			limit, err = strconv.Atoi(limitStr)
			if err != nil || limit <= 0 {
				JSONError(w, "Limite (limit) deve ser um número positivo", http.StatusBadRequest)
				return
			}
		}

		offsetStr := r.URL.Query().Get("offset")
		var offset int
		if offsetStr != "" {
			var err error
			offset, err = strconv.Atoi(offsetStr)
			if err != nil || offset < 0 {
				JSONError(w, "Deslocamento (offset) deve ser um número não negativo", http.StatusBadRequest)
				return
			}
		}
		sortParam := r.URL.Query().Get("sort")

		var materiais interface{}
		var err error

		if termo != "" || categoria != "" || fonte != "" || anoInicio > 0 || anoFim > 0 {
			materiais, err = pesquisarUC.Execute(r.Context(), termo, categoria, fonte, anoInicio, anoFim, nil, limit, offset, sortParam)
		} else {
			materiais, err = listarUC.Execute(r.Context(), limit, offset)
		}

		if err != nil {
			JSONError(w, "Erro ao buscar materiais: "+err.Error(), http.StatusInternalServerError)
			return
		}

		JSONSuccess(w, materiais, http.StatusOK)
	})

	mux.HandleFunc("GET /materiais/detalhes", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			JSONError(w, "ID inválido", http.StatusBadRequest)
			return
		}

		m, err := buscarUC.Execute(r.Context(), id)
		if err != nil {
			JSONError(w, "Material não encontrado", http.StatusNotFound)
			return
		}

		JSONSuccess(w, m, http.StatusOK)
	})

	mux.HandleFunc("GET /materiais/similares", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			JSONError(w, "ID inválido", http.StatusBadRequest)
			return
		}

		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 4
		}

		materiais, err := similaresUC.Execute(r.Context(), id, limit)
		if err != nil {
			JSONError(w, "Erro ao buscar materiais similares: "+err.Error(), http.StatusInternalServerError)
			return
		}

		JSONSuccess(w, materiais, http.StatusOK)
	})

	mux.HandleFunc("POST /materiais/favoritar", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UsuarioID  int  `json:"usuario_id"`
			MaterialID int  `json:"material_id"`
			Favoritar  bool `json:"favoritar"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if err := favoritarUC.Execute(r.Context(), req.UsuarioID, req.MaterialID, req.Favoritar); err != nil {
			JSONError(w, "Erro ao favoritar material: "+err.Error(), http.StatusInternalServerError)
			return
		}

		JSONSuccess(w, nil, http.StatusOK)
	})

	mux.HandleFunc("GET /materiais/favoritos", func(w http.ResponseWriter, r *http.Request) {
		uIDStr := r.URL.Query().Get("usuario_id")
		usuarioID, err := strconv.Atoi(uIDStr)
		if err != nil || usuarioID <= 0 {
			JSONError(w, "ID de usuário inválido", http.StatusBadRequest)
			return
		}
		favoritos, err := favoritarUC.Listar(r.Context(), usuarioID)
		if err != nil {
			JSONError(w, "Erro ao listar favoritos", http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, favoritos, http.StatusOK)
	})

	mux.HandleFunc("POST /materiais/avaliar", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UsuarioID  int     `json:"usuario_id"`
			MaterialID int     `json:"material_id"`
			Nota       float64 `json:"nota"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if err := avaliarUC.Execute(r.Context(), req.UsuarioID, req.MaterialID, req.Nota); err != nil {
			JSONError(w, "Erro ao avaliar material: "+err.Error(), http.StatusInternalServerError)
			return
		}

		JSONSuccess(w, nil, http.StatusOK)
	})

	mux.HandleFunc("/materiais/historico", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req struct {
				UsuarioID  int `json:"usuario_id"`
				MaterialID int `json:"material_id"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				JSONError(w, "JSON inválido", http.StatusBadRequest)
				return
			}
			if err := historicoUC.Execute(r.Context(), req.UsuarioID, req.MaterialID); err != nil {
				JSONError(w, "Erro ao registrar histórico: "+err.Error(), http.StatusInternalServerError)
				return
			}
			JSONSuccess(w, nil, http.StatusOK)
		case http.MethodGet:
			uIDStr := r.URL.Query().Get("usuario_id")
			usuarioID, err := strconv.Atoi(uIDStr)
			if err != nil || usuarioID <= 0 {
				JSONError(w, "ID de usuário inválido", http.StatusBadRequest)
				return
			}
			historico, err := historicoUC.Listar(r.Context(), usuarioID)
			if err != nil {
				JSONError(w, "Erro ao listar histórico: "+err.Error(), http.StatusInternalServerError)
				return
			}
			JSONSuccess(w, historico, http.StatusOK)
		default:
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("GET /materiais/pdf/proxy", func(w http.ResponseWriter, r *http.Request) {
		pdfURL := r.URL.Query().Get("url")
		if pdfURL == "" {
			JSONError(w, "URL do PDF é obrigatória", http.StatusBadRequest)
			return
		}

		// Segurança Básica SSRF: validar apenas HTTP/HTTPS e bloquear IPs locais/privados
		parsedURL, err := url.Parse(pdfURL)
		if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
			JSONError(w, "URL inválida - apenas HTTP/HTTPS são permitidos", http.StatusBadRequest)
			return
		}

		// Resolução de IP e verificação de SSRF profunda
		ips, err := net.LookupIP(parsedURL.Hostname())
		if err != nil || len(ips) == 0 {
			JSONError(w, "Não foi possível resolver o host", http.StatusBadRequest)
			return
		}
		for _, ip := range ips {
			if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
				JSONError(w, "Acesso bloqueado a IPs internos ou privados por segurança (SSRF Protection)", http.StatusForbidden)
				return
			}
		}

		// Cliente HTTP seguro com Timeout
		safeClient := &http.Client{
			Timeout: 15 * time.Second,
		}

		// Faz o fetch do PDF original
		resp, err := safeClient.Get(pdfURL)
		if err != nil {
			JSONError(w, "Erro ao buscar PDF original: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			JSONError(w, fmt.Sprintf("Erro na fonte original: %d %s", resp.StatusCode, resp.Status), http.StatusBadGateway)
			return
		}

		// Garantir que é um PDF
		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/pdf" {
			JSONError(w, "O recurso solicitado não é um PDF válido", http.StatusBadRequest)
			return
		}

		// Copia os headers relevantes para o browser
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline; filename=\"document.pdf\"")

		// Prevenção de XSS e Clickjacking
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'self';")

		// Proteção contra arquivos gigantes (DoS): Limita a 50 MB
		limitReader := io.LimitReader(resp.Body, 50<<20)

		if _, err = io.Copy(w, limitReader); err != nil {
			// Não use JSONError se o body já tiver começado a ser enviado
			return
		}
	})

}
