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

	mux.HandleFunc("/materiais", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(materiais)
	})

	mux.HandleFunc("/materiais/detalhes", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		m, err := buscarUC.Execute(id)
		if err != nil {
			http.Error(w, "Material não encontrado", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(m)
	})

	mux.HandleFunc("/materiais/recomendacoes", func(w http.ResponseWriter, r *http.Request) {
		usuarioID, _ := strconv.Atoi(r.URL.Query().Get("usuario_id"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 5
		}

		materiais, err := recomendacaoUC.Execute(usuarioID, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(materiais)
	})
}
