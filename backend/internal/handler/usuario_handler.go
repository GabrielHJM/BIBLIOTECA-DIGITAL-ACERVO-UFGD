package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"biblioteca-digital-api/internal/dto"
	"biblioteca-digital-api/internal/usecase/usuario"
	domain "biblioteca-digital-api/internal/domain/usuario"
)

func RegisterUsuarioRoutes(mux *http.ServeMux, db *sql.DB) {
	// Aqui usaríamos um repository real, simulando agora:
	repo := NewUsuarioPG(db)
	cadastrarUC := usuario.NewCadastrarUsuario(repo)
	loginUC := usuario.NewLoginUseCase(repo)

	mux.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		var req dto.UsuarioRequest
		json.NewDecoder(r.Body).Decode(&req)
		u := domain.Usuario{Nome: req.Nome, Email: req.Email, Senha: req.Senha, Tipo: req.Tipo}
		err := cadastrarUC.Execute(r.Context(), u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}
		var req dto.LoginRequest
		json.NewDecoder(r.Body).Decode(&req)
		token, err := loginUC.Execute(r.Context(), req.Email, req.Senha)
		if err != nil {
			http.Error(w, "Login inválido", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	})
}
