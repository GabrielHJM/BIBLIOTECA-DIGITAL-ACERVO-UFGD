package handler

import (
	"biblioteca-digital-api/internal/repository"
	"biblioteca-digital-api/internal/usecase/material"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func RegisterMaterialRoutes(mux *http.ServeMux, db *sql.DB) {
	repo := &repository.MaterialPostgres{DB: db}
	listarUC := &material.ListarConteudosUseCase{Repo: repo}
	buscarUC := &material.BuscarMaterialUseCase{Repo: repo}
	pesquisarUC := &material.PesquisarMaterialUseCase{Repo: repo}
	recomendacaoUC := &material.ObterRecomendacoesUseCase{Repo: repo}
	favoritarUC := &material.FavoritarMaterialUseCase{Repo: repo}
	avaliarUC := &material.AvaliarMaterialUseCase{Repo: repo}

	mux.HandleFunc("/materiais", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		termo := r.URL.Query().Get("q")
		categoria := r.URL.Query().Get("categoria")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 10
		}
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

		var materiais interface{}
		var err error

		if termo != "" || categoria != "" {
			materiais, err = pesquisarUC.Execute(termo, categoria, nil, limit, offset)
		} else {
			materiais, err = listarUC.Execute(limit, offset)
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

		m, err := buscarUC.Execute(id)
		if err != nil {
			JSONError(w, "Material não encontrado", http.StatusNotFound)
			return
		}

		JSONSuccess(w, m, http.StatusOK)
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

		materiais, err := recomendacaoUC.Execute(usuarioID, limit)
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

		if err := favoritarUC.Execute(req.UsuarioID, req.MaterialID, req.Favoritar); err != nil {
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
		favoritos, err := favoritarUC.Listar(usuarioID)
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

		if err := avaliarUC.Execute(req.UsuarioID, req.MaterialID, req.Nota, req.Comentario); err != nil {
			JSONError(w, "Erro ao avaliar material: "+err.Error(), http.StatusInternalServerError)
			return
		}

		JSONSuccess(w, nil, http.StatusOK)
	})
}
