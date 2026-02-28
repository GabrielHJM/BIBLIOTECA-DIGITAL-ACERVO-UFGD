package handler

import (
	"biblioteca-digital-api/internal/repository"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func RegisterAdminRoutes(mux *http.ServeMux, db *sql.DB) {
	repo := &repository.MaterialPostgres{DB: db}

	// Listar Materiais Pendentes
	mux.HandleFunc("GET /admin/materiais/pendentes", func(w http.ResponseWriter, r *http.Request) {
		materiais, err := repo.ListarPendentes(r.Context())
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, materiais, http.StatusOK)
	})

	// Aprovar/Rejeitar Material
	mux.HandleFunc("POST /admin/materiais/moderar", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID        int    `json:"id"`
			Status    string `json:"status"` // 'aprovado' ou 'rejeitado'
			CuradorID int    `json:"curador_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if err := repo.AtualizarStatus(r.Context(), req.ID, req.Status, req.CuradorID); err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, nil, http.StatusOK)
	})

	// Métricas Globais
	mux.HandleFunc("GET /admin/metricas", func(w http.ResponseWriter, r *http.Request) {
		metricas, err := repo.ObterMetricasGlobais(r.Context())
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, metricas, http.StatusOK)
	})

	// Deletar material
	mux.HandleFunc("DELETE /admin/materiais/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		if err := repo.Deletar(r.Context(), id); err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, nil, http.StatusOK)
	})
}
