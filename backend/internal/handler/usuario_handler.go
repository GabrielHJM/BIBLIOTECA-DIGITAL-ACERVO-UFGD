package handler

import (
	domain "biblioteca-digital-api/internal/domain/usuario"
	"biblioteca-digital-api/internal/dto"
	"biblioteca-digital-api/internal/repository"
	"biblioteca-digital-api/internal/usecase/usuario"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func RegisterUsuarioRoutes(mux *http.ServeMux, db *sql.DB) {
	repo := repository.NewUsuarioPG(db)
	cadastrarUC := usuario.NewCadastrarUsuario(repo)
	loginUC := usuario.NewLoginUseCase(repo)
	redefinirSenhaUC := usuario.NewRedefinirSenhaUseCase(repo)
	atualizarUC := usuario.NewAtualizarUsuario(repo)

	mux.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		var req dto.UsuarioRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		u := domain.Usuario{Nome: req.Nome, Email: req.Email, Senha: req.Senha, Tipo: req.Tipo}
		err := cadastrarUC.Execute(r.Context(), u)
		if err != nil {
			JSONError(w, "Erro ao cadastrar usuário: "+err.Error(), http.StatusUnprocessableEntity)
			return
		}
		JSONSuccess(w, nil, http.StatusCreated)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		var req dto.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Login: Erro ao decodificar JSON")
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		log.Printf("Tentativa de login para: %s", req.Email)

		token, err := loginUC.Execute(r.Context(), req.Email, req.Senha)
		if err != nil {
			log.Printf("Login falhou para %s: %v", req.Email, err)
			JSONError(w, "Login inválido: credenciais incorretas", http.StatusUnauthorized)
			return
		}

		// Buscar dados do usuário para retornar ao frontend
		u, _ := repo.BuscarPorEmail(r.Context(), req.Email)

		log.Printf("Login bem-sucedido: %s (ID: %d)", req.Email, u.ID)

		JSONSuccess(w, map[string]interface{}{
			"token":    token,
			"id":       u.ID,
			"nome":     u.Nome,
			"email":    u.Email,
			"foto_url": u.FotoURL,
		}, http.StatusOK)
	})

	mux.HandleFunc("/redefinir-senha", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		var req dto.RedefinirSenhaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		err := redefinirSenhaUC.Execute(r.Context(), req.Email, req.Senha)
		if err != nil {
			JSONError(w, "Erro ao redefinir senha: "+err.Error(), http.StatusUnprocessableEntity)
			return
		}
		JSONSuccess(w, nil, http.StatusOK)
	})

	mux.HandleFunc("PUT /usuarios/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		log.Printf("Recebido pedido de atualização para ID: %s", idStr)
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Erro: ID inválido %s", idStr)
			JSONError(w, "ID inválido", http.StatusBadRequest)
			return
		}

		var req dto.UsuarioRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Erro: JSON inválido")
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		log.Printf("Dados recebidos: Nome=%s, Email=%s", req.Nome, req.Email)

		u := domain.Usuario{
			ID:      id,
			Nome:    req.Nome,
			Email:   req.Email,
			FotoURL: req.FotoURL,
		}

		err = atualizarUC.Execute(r.Context(), u)
		if err != nil {
			log.Printf("Erro ao atualizar no banco: %v", err)
			JSONError(w, "Erro ao atualizar usuário: "+err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Usuário %d atualizado com sucesso", id)
		JSONSuccess(w, nil, http.StatusOK)
	})
}
