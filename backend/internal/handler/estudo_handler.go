package handler

import (
	"biblioteca-digital-api/internal/domain/estudo"
	"biblioteca-digital-api/internal/repository"
	usecase "biblioteca-digital-api/internal/usecase/estudo"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func RegisterEstudoRoutes(mux *http.ServeMux, db *sql.DB) {
	repo := repository.NewEstudoPostgres(db)
	matRepo := &repository.MaterialPostgres{DB: db}
	uc := usecase.NewUseCase(repo, matRepo)

	// Flashcards
	mux.HandleFunc("GET /estudo/flashcards", func(w http.ResponseWriter, r *http.Request) {
		uID, err := strconv.Atoi(r.URL.Query().Get("usuario_id"))
		if err != nil || uID <= 0 {
			JSONError(w, "ID de usuário inválido", http.StatusBadRequest)
			return
		}
		mID, _ := strconv.Atoi(r.URL.Query().Get("material_id"))
		cards, err := uc.ListarFlashcards(r.Context(), uID, mID)
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

	mux.HandleFunc("POST /estudo/flashcards", func(w http.ResponseWriter, r *http.Request) {
		var f estudo.Flashcard
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		if err := uc.CriarFlashcard(r.Context(), &f); err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, f, http.StatusCreated)
	})

	mux.HandleFunc("DELETE /estudo/flashcards", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id <= 0 {
			JSONError(w, "ID do flashcard inválido", http.StatusBadRequest)
			return
		}
		uID, err := strconv.Atoi(r.URL.Query().Get("usuario_id"))
		if err != nil || uID <= 0 {
			JSONError(w, "ID de usuário inválido", http.StatusBadRequest)
			return
		}
		if err := repo.DeletarFlashcard(r.Context(), id, uID); err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, nil, http.StatusOK)
	})
}
