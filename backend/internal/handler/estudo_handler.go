package handler

import (
	"biblioteca-digital-api/internal/domain/estudo"
	"biblioteca-digital-api/internal/pkg/ai"
	"biblioteca-digital-api/internal/repository"
	usecase "biblioteca-digital-api/internal/usecase/estudo"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func RegisterEstudoRoutes(mux *http.ServeMux, db *sql.DB, gemini *ai.GeminiClient) {
	repo := repository.NewEstudoPostgres(db)
	matRepo := &repository.MaterialPostgres{DB: db}
	uc := usecase.NewUseCase(repo, matRepo, gemini)

	// Anotações
	mux.HandleFunc("POST /estudo/anotacoes", func(w http.ResponseWriter, r *http.Request) {
		var a estudo.Anotacao
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		if err := uc.CriarAnotacao(r.Context(), &a); err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, a, http.StatusCreated)
	})

	mux.HandleFunc("GET /estudo/anotacoes", func(w http.ResponseWriter, r *http.Request) {
		uID, _ := strconv.Atoi(r.URL.Query().Get("usuario_id"))
		mID, _ := strconv.Atoi(r.URL.Query().Get("material_id"))
		anotacoes, err := uc.ListarAnotacoes(r.Context(), uID, mID)
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, anotacoes, http.StatusOK)
	})

	mux.HandleFunc("DELETE /estudo/anotacoes", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		uID, _ := strconv.Atoi(r.URL.Query().Get("usuario_id"))
		if err := uc.DeletarAnotacao(r.Context(), id, uID); err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, nil, http.StatusOK)
	})

	// Flashcards
	mux.HandleFunc("GET /estudo/flashcards", func(w http.ResponseWriter, r *http.Request) {
		uID, _ := strconv.Atoi(r.URL.Query().Get("usuario_id"))
		mID, _ := strconv.Atoi(r.URL.Query().Get("material_id"))
		cards, err := uc.ListarFlashcards(r.Context(), uID, mID)
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, cards, http.StatusOK)
	})

	mux.HandleFunc("POST /estudo/flashcards/gerar", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UsuarioID  int `json:"usuario_id"`
			MaterialID int `json:"material_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		cards, err := uc.GerarFlashcardsIA(r.Context(), req.UsuarioID, req.MaterialID)
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, cards, http.StatusOK)
	})

	mux.HandleFunc("PUT /estudo/flashcards/revisar", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID          int `json:"id"`
			Dificuldade int `json:"dificuldade"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		if err := uc.AtualizarDificuldade(r.Context(), req.ID, req.Dificuldade); err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, nil, http.StatusOK)
	})
}
