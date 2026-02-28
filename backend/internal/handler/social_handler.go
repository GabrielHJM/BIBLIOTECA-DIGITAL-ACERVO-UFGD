package handler

import (
	"biblioteca-digital-api/internal/domain/social"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// RegisterSocialRoutes registers HTTP endpoints for social interactions.
func RegisterSocialRoutes(mux *http.ServeMux, uc social.UseCase, repo social.ComentarioRepository, likeRepo social.LikeRepository) {

	mux.HandleFunc("POST /material/like", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			UsuarioID  int `json:"usuario_id"`
			MaterialID int `json:"material_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		liked, err := uc.ToggleLike(payload.UsuarioID, payload.MaterialID)
		if err != nil {
			JSONError(w, "Erro ao processar like", http.StatusInternalServerError)
			return
		}

		JSONSuccess(w, map[string]interface{}{"liked": liked}, http.StatusOK)
	})

	mux.HandleFunc("POST /material/comment", func(w http.ResponseWriter, r *http.Request) {
		var c social.Comentario
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		c.CreatedAt = time.Now()
		if err := uc.AddComment(c); err != nil {
			JSONError(w, "Erro ao salvar comentário", http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, nil, http.StatusCreated)
	})

	mux.HandleFunc("POST /material/share", func(w http.ResponseWriter, r *http.Request) {
		var s social.ShareLog
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		s.SharedAt = time.Now()
		if err := uc.RegisterShare(s); err != nil {
			JSONError(w, "Erro ao registrar share", http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, nil, http.StatusCreated)
	})

	mux.HandleFunc("POST /material/message", func(w http.ResponseWriter, r *http.Request) {
		var m social.Mensagem
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		m.SentAt = time.Now()
		if err := uc.SendMessage(m); err != nil {
			JSONError(w, "Erro ao enviar mensagem", http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, nil, http.StatusCreated)
	})

	// Get likes info
	mux.HandleFunc("GET /material/likes", func(w http.ResponseWriter, r *http.Request) {
		materialIDStr := r.URL.Query().Get("material_id")
		materialID, err := strconv.Atoi(materialIDStr)
		if err != nil {
			JSONError(w, "material_id inválido", http.StatusBadRequest)
			return
		}
		// Optional user_id to indicate if user liked
		userIDStr := r.URL.Query().Get("usuario_id")
		var likedByUser bool
		if userIDStr != "" {
			uid, err := strconv.Atoi(userIDStr)
			if err == nil {
				likedByUser, _ = likeRepo.HasLiked(uid, materialID)
			}
		}
		count, err := likeRepo.CountLikes(materialID)
		if err != nil {
			JSONError(w, "Erro ao contar likes", http.StatusInternalServerError)
			return
		}
		resp := map[string]interface{}{"count": count, "likedByUser": likedByUser}
		JSONSuccess(w, resp, http.StatusOK)
	})

	// Get comments list
	mux.HandleFunc("GET /material/comments", func(w http.ResponseWriter, r *http.Request) {
		materialIDStr := r.URL.Query().Get("material_id")
		materialID, err := strconv.Atoi(materialIDStr)
		if err != nil {
			JSONError(w, "material_id inválido", http.StatusBadRequest)
			return
		}
		comments, err := repo.ListComentarios(materialID)
		if err != nil {
			JSONError(w, "Erro ao listar comentários", http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, comments, http.StatusOK)
	})

}
