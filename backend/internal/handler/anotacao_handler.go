package handler

import (
	domain "biblioteca-digital-api/internal/domain/anotacao"
	"biblioteca-digital-api/internal/handler/middleware"
	"biblioteca-digital-api/internal/pkg/logger"
	"biblioteca-digital-api/internal/repository"
	anotacaoUsecase "biblioteca-digital-api/internal/usecase/anotacao"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

func RegisterAnotacaoRoutes(mux *http.ServeMux, db *sql.DB) {
	repo := repository.NewAnotacaoRepositoryPG(db)
	uc := anotacaoUsecase.NewUsecase(repo)

	// Helper para ler usuario_id do contexto com o tipo correto do middleware
	getUsuarioID := func(r *http.Request) (int, bool) {
		v := r.Context().Value(middleware.UsuarioIDKey)
		if v == nil {
			return 0, false
		}
		id, ok := v.(int)
		return id, ok
	}

	mux.Handle("POST /anotacoes", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usuarioID, ok := getUsuarioID(r)
		if !ok {
			JSONError(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		var req domain.Anotacao
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		req.UsuarioID = usuarioID

		id, err := uc.Criar(r.Context(), req)
		if err != nil {
			logger.Error("Erro ao criar anotação", zap.Error(err))
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		JSONSuccess(w, map[string]interface{}{"id": id, "mensagem": "Anotação criada"}, http.StatusCreated)
	})))

	mux.Handle("GET /anotacoes", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usuarioID, ok := getUsuarioID(r)
		if !ok {
			JSONError(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		lista, err := uc.ListarPorUsuario(r.Context(), usuarioID)
		if err != nil {
			logger.Error("Erro ao listar anotações", zap.Error(err))
			JSONError(w, "Erro ao buscar anotações", http.StatusInternalServerError)
			return
		}
		if lista == nil {
			lista = []domain.Anotacao{}
		}
		JSONSuccess(w, lista, http.StatusOK)
	})))

	mux.Handle("PUT /anotacoes/{id}", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usuarioID, ok := getUsuarioID(r)
		if !ok {
			JSONError(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			JSONError(w, "ID inválido", http.StatusBadRequest)
			return
		}

		var req domain.Anotacao
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		req.ID = id
		req.UsuarioID = usuarioID

		if err := uc.Atualizar(r.Context(), req); err != nil {
			logger.Error("Erro ao atualizar anotação", zap.Error(err))
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, map[string]interface{}{"mensagem": "Anotação atualizada"}, http.StatusOK)
	})))

	mux.Handle("DELETE /anotacoes/{id}", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usuarioID, ok := getUsuarioID(r)
		if !ok {
			JSONError(w, "Usuário não autenticado", http.StatusUnauthorized)
			return
		}

		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			JSONError(w, "ID inválido", http.StatusBadRequest)
			return
		}

		if err := uc.Excluir(r.Context(), id, usuarioID); err != nil {
			logger.Error("Erro ao excluir anotação", zap.Error(err))
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, map[string]interface{}{"mensagem": "Anotação excluída"}, http.StatusOK)
	})))
}
