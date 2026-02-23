package handler

import (
	domain "biblioteca-digital-api/internal/domain/usuario"
	"biblioteca-digital-api/internal/dto"
	"biblioteca-digital-api/internal/repository"
	"biblioteca-digital-api/internal/usecase/usuario"
	"database/sql"
	"encoding/json"
	"net/http"
)

func RegisterUsuarioRoutes(mux *http.ServeMux, db *sql.DB) {
	repo := repository.NewUsuarioPG(db)
	cadastrarUC := usuario.NewCadastrarUsuario(repo)
	loginUC := usuario.NewLoginUseCase(repo)
	redefinirSenhaUC := usuario.NewRedefinirSenhaUseCase(repo)

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
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		token, err := loginUC.Execute(r.Context(), req.Email, req.Senha)
		if err != nil {
			JSONError(w, "Login inválido: credenciais incorretas", http.StatusUnauthorized)
			return
		}
		JSONSuccess(w, map[string]string{"token": token}, http.StatusOK)
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
}
