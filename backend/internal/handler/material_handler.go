package handler

import (
	"biblioteca-digital-api/internal/harvester"
	"biblioteca-digital-api/internal/pkg/ai"
	"biblioteca-digital-api/internal/pkg/cache"
	"biblioteca-digital-api/internal/pkg/metadata"
	"biblioteca-digital-api/internal/repository"
	"biblioteca-digital-api/internal/usecase/material"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func RegisterMaterialRoutes(mux *http.ServeMux, db *sql.DB, gemini *ai.GeminiClient, c cache.Cache) {
	repo := &repository.MaterialPostgres{DB: db}
	mh := harvester.NewMultiSourceHarvester()
	meta := metadata.NewMetadataService()

	listarUC := &material.ListarConteudosUseCase{Repo: repo, Harvester: mh, Cache: c}
	buscarUC := &material.BuscarMaterialUseCase{Repo: repo, Meta: meta}
	similaresUC := &material.BuscarSimilaresUseCase{Repo: repo}
	pesquisarUC := &material.PesquisarMaterialUseCase{Repo: repo, Harvester: mh, Cache: c}
	recomendacaoUC := &material.ObterRecomendacoesUseCase{Repo: repo}
	favoritarUC := &material.FavoritarMaterialUseCase{Repo: repo}
	avaliarUC := &material.AvaliarMaterialUseCase{Repo: repo}
	emprestarUC := &material.CriarEmprestimoUseCase{Repo: repo}
	historicoUC := &material.HistoricoLeituraUseCase{Repo: repo}
	aiUC := material.NewAIUseCase(repo, gemini)

	mux.HandleFunc("GET /materiais", func(w http.ResponseWriter, r *http.Request) {
		termo := r.URL.Query().Get("q")
		categoria := r.URL.Query().Get("categoria")
		fonte := r.URL.Query().Get("fonte")
		anoInicio, _ := strconv.Atoi(r.URL.Query().Get("ano_inicio"))
		anoFim, _ := strconv.Atoi(r.URL.Query().Get("ano_fim"))

		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 10
		}
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
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

	mux.HandleFunc("GET /materiais/recomendacoes", func(w http.ResponseWriter, r *http.Request) {
		usuarioID, _ := strconv.Atoi(r.URL.Query().Get("usuario_id"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 5
		}

		materiais, err := recomendacaoUC.Execute(r.Context(), usuarioID, limit)
		if err != nil {
			JSONError(w, "Erro ao obter recomendações: "+err.Error(), http.StatusInternalServerError)
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
		usuarioID, _ := strconv.Atoi(r.URL.Query().Get("usuario_id"))
		favoritos, err := favoritarUC.Listar(r.Context(), usuarioID)
		if err != nil {
			JSONError(w, "Erro ao listar favoritos", http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, favoritos, http.StatusOK)
	})

	mux.HandleFunc("POST /materiais/avaliar", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UsuarioID  int    `json:"usuario_id"`
			MaterialID int    `json:"material_id"`
			Nota       int    `json:"nota"`
			Comentario string `json:"comentario"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if err := avaliarUC.Execute(r.Context(), req.UsuarioID, req.MaterialID, req.Nota, req.Comentario); err != nil {
			JSONError(w, "Erro ao avaliar material: "+err.Error(), http.StatusInternalServerError)
			return
		}

		JSONSuccess(w, nil, http.StatusOK)
	})

	mux.HandleFunc("POST /materiais/emprestar", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UsuarioID  int `json:"usuario_id"`
			MaterialID int `json:"material_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		if err := emprestarUC.Execute(r.Context(), req.UsuarioID, req.MaterialID); err != nil {
			JSONError(w, "Erro ao realizar empréstimo: "+err.Error(), http.StatusInternalServerError)
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
			usuarioID, _ := strconv.Atoi(r.URL.Query().Get("usuario_id"))
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

	mux.HandleFunc("GET /materiais/avaliacoes", func(w http.ResponseWriter, r *http.Request) {
		materialID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			JSONError(w, "ID do material inválido", http.StatusBadRequest)
			return
		}

		avaliacoes, err := avaliarUC.Listar(r.Context(), materialID)
		if err != nil {
			JSONError(w, "Erro ao listar avaliações", http.StatusInternalServerError)
			return
		}

		JSONSuccess(w, avaliacoes, http.StatusOK)
	})

	mux.HandleFunc("POST /materiais/{id}/chat", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			JSONError(w, "ID inválido", http.StatusBadRequest)
			return
		}

		var req struct {
			Pergunta string `json:"pergunta"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		resposta, err := aiUC.PerguntarLivro(r.Context(), id, req.Pergunta)
		if err != nil {
			JSONError(w, "Erro na IA: "+err.Error(), http.StatusInternalServerError)
			return
		}

		JSONSuccess(w, map[string]string{"resposta": resposta}, http.StatusOK)
	})

	mux.HandleFunc("GET /materiais/{id}/resumo", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			JSONError(w, "ID inválido", http.StatusBadRequest)
			return
		}

		resumo, err := aiUC.GerarResumo(r.Context(), id)
		if err != nil {
			JSONError(w, "Erro ao gerar resumo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		JSONSuccess(w, map[string]string{"resumo": resumo}, http.StatusOK)
	})

	mux.HandleFunc("GET /materiais/pdf/proxy", func(w http.ResponseWriter, r *http.Request) {
		pdfURL := r.URL.Query().Get("url")
		if pdfURL == "" {
			JSONError(w, "URL do PDF é obrigatória", http.StatusBadRequest)
			return
		}

		// Segurança Básica SSRF: validar apenas HTTP/HTTPS
		if len(pdfURL) < 7 || (pdfURL[:7] != "http://" && pdfURL[:8] != "https://") {
			JSONError(w, "URL inválida - apenas HTTP/HTTPS são permitidos", http.StatusBadRequest)
			return
		}

		// Faz o fetch do PDF original
		resp, err := http.Get(pdfURL)
		if err != nil {
			JSONError(w, "Erro ao buscar PDF original: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			JSONError(w, "Erro na fonte original: "+resp.Status, http.StatusBadGateway)
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

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			// Aqui não podemos mais usar JSONError pois o body pode já ter começado a carregar
			return
		}
	})

}
