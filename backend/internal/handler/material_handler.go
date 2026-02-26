package handler

import (
	"biblioteca-digital-api/internal/harvester"
	"biblioteca-digital-api/internal/pkg/cache"
	"biblioteca-digital-api/internal/repository"
	"biblioteca-digital-api/internal/usecase/material"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func RegisterMaterialRoutes(mux *http.ServeMux, db *sql.DB) {
	repo := &repository.MaterialPostgres{DB: db}
	mh := harvester.NewMultiSourceHarvester()
	c := cache.NewMemoryCache()

	listarUC := &material.ListarConteudosUseCase{Repo: repo, Harvester: mh, Cache: c}
	buscarUC := &material.BuscarMaterialUseCase{Repo: repo}
	similaresUC := &material.BuscarSimilaresUseCase{Repo: repo}
	pesquisarUC := &material.PesquisarMaterialUseCase{Repo: repo, Harvester: mh, Cache: c}
	recomendacaoUC := &material.ObterRecomendacoesUseCase{Repo: repo}
	favoritarUC := &material.FavoritarMaterialUseCase{Repo: repo}
	avaliarUC := &material.AvaliarMaterialUseCase{Repo: repo}
	emprestarUC := &material.CriarEmprestimoUseCase{Repo: repo}
	historicoUC := &material.HistoricoLeituraUseCase{Repo: repo}

	mux.HandleFunc("/materiais", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

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
			// Nota: O Repository ainda não suporta fonte/ano na busca local (FTS),
			// mas o Harvester pode usar. Implementaremos o filtro no Repo se necessário.
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

	mux.HandleFunc("/materiais/detalhes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

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

	mux.HandleFunc("/materiais/similares", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

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

	mux.HandleFunc("/materiais/recomendacoes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

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

	mux.HandleFunc("/materiais/favoritar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

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

	mux.HandleFunc("/materiais/favoritos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		usuarioID, _ := strconv.Atoi(r.URL.Query().Get("usuario_id"))
		favoritos, err := favoritarUC.Listar(r.Context(), usuarioID)
		if err != nil {
			JSONError(w, "Erro ao listar favoritos", http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, favoritos, http.StatusOK)
	})

	mux.HandleFunc("/materiais/avaliar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

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

	mux.HandleFunc("/materiais/emprestar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
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

	mux.HandleFunc("/materiais/avaliacoes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

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
}
