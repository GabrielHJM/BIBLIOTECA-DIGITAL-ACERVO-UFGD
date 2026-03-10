package handler

import (
	"biblioteca-digital-api/internal/domain/notificacao"
	"biblioteca-digital-api/internal/repository"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func RegisterNotificacaoRoutes(mux *http.ServeMux, db *sql.DB) {
	repo := repository.NewNotificacaoPostgres(db)

	mux.HandleFunc("GET /notificacoes", func(w http.ResponseWriter, r *http.Request) {
		uID, _ := strconv.Atoi(r.URL.Query().Get("usuario_id"))
		ns, err := repo.ListarPorUsuario(r.Context(), uID)
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, ns, http.StatusOK)
	})

	mux.HandleFunc("POST /notificacoes", func(w http.ResponseWriter, r *http.Request) {
		var n notificacao.Notificacao
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		if err := repo.Criar(r.Context(), &n); err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, n, http.StatusCreated)
	})

	mux.HandleFunc("PUT /notificacoes/ler", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		if err := repo.MarcarComoLida(r.Context(), id); err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, nil, http.StatusOK)
	})

	mux.HandleFunc("DELETE /notificacoes", func(w http.ResponseWriter, r *http.Request) {
		uID, _ := strconv.Atoi(r.URL.Query().Get("usuario_id"))
		if err := repo.LimparPorUsuario(r.Context(), uID); err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, nil, http.StatusOK)
	})
}
