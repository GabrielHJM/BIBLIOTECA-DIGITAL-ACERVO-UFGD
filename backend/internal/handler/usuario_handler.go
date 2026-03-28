package handler

import (
	domain "biblioteca-digital-api/internal/domain/usuario"
	"biblioteca-digital-api/internal/dto"
	"biblioteca-digital-api/internal/pkg/logger"
	"biblioteca-digital-api/internal/repository"
	"biblioteca-digital-api/internal/usecase/usuario"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"biblioteca-digital-api/internal/pkg/validation"
	"biblioteca-digital-api/pkg/auth"
	"biblioteca-digital-api/pkg/hash"
	"time"

	"go.uber.org/zap"
)

func RegisterUsuarioRoutes(mux *http.ServeMux, db *sql.DB) {
	repo := repository.NewUsuarioPG(db)
	cadastrarUC := usuario.NewCadastrarUsuario(repo)
	loginUC := usuario.NewLoginUseCase(repo)
	redefinirSenhaUC := usuario.NewRedefinirSenhaUseCase(repo)
	atualizarUC := usuario.NewAtualizarUsuario(repo)
	atualizarMetaUC := usuario.NewAtualizarMeta(repo)

	// @Summary Cadastrar usuário
	// @Description Cria um novo usuário no sistema
	// @Tags usuários
	// @Accept json
	// @Produce json
	// @Param usuário body dto.UsuarioRequest true "Dados do usuário"
	// @Success 201 {object} Response
	// @Failure 400 {object} Response
	// @Router /usuarios [post]
	mux.HandleFunc("POST /usuarios", func(w http.ResponseWriter, r *http.Request) {
		var req dto.UsuarioRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		u := domain.Usuario{
			Nome:     req.Nome,
			Email:    req.Email,
			Senha:    req.Senha,
			Tipo:     req.Tipo,
			Username: req.Username,
			Cpf:      req.Cpf,
		}

		if req.DataNascimento != "" {
			t, err := time.Parse("2006-01-02", req.DataNascimento)
			if err == nil {
				u.DataNascimento = &t
			}
		}

		err := cadastrarUC.Execute(r.Context(), &u)
		if err != nil {
			logger.Error("Erro ao cadastrar usuário", zap.Error(err))
			JSONError(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		logger.Info("Usuário cadastrado com sucesso", zap.Int("id", u.ID), zap.String("email", u.Email))

		// Auto-login: Gerar token após cadastro bem sucedido
		token, err := auth.GerarToken(u.ID)
		if err != nil {
			logger.Error("Erro ao gerar token pós-cadastro", zap.Error(err))
			JSONSuccess(w, nil, http.StatusCreated)
			return
		}

		JSONSuccess(w, map[string]interface{}{
			"token":    token,
			"id":       u.ID,
			"nome":     u.Nome,
			"email":    u.Email,
			"username": u.Username,
			"foto_url": u.FotoURL,
		}, http.StatusCreated)
	})

	// @Summary Login de usuário
	// @Description Autentica um usuário e retorna um token JWT
	// @Tags usuários
	// @Accept json
	// @Produce json
	// @Param login body dto.LoginRequest true "Credenciais de login"
	// @Success 200 {object} Response{data=map[string]interface{}}
	// @Failure 401 {object} Response
	// @Router /login [post]
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		var req dto.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Login: Erro ao decodificar JSON", zap.Error(err))
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if err := validation.ValidateStruct(req); err != nil {
			logger.Warn("Login: Dados inválidos", zap.Error(err))
			JSONError(w, "Dados de login inválidos", http.StatusBadRequest)
			return
		}

		logger.Info("Tentativa de login", zap.String("email", req.Email))

		// O loginUC.Execute já deve validar a senha. Idealmente retornaria o usuário também.
		token, err := loginUC.Execute(r.Context(), req.Email, req.Senha)

		// Otimização: Buscar dados apenas se o login UC não retornar.
		// Se o Repository suportar busca por email (como já suporta), usamos aqui.
		u, err := repo.BuscarPorEmail(r.Context(), req.Email)
		if err != nil {
			logger.Error("Erro ao recuperar dados do usuário após login", zap.String("email", req.Email), zap.Error(err))
			JSONError(w, "Erro interno ao processar login", http.StatusInternalServerError)
			return
		}

		logger.Info("Login bem-sucedido", zap.String("email", req.Email), zap.Int("id", u.ID))

		JSONSuccess(w, map[string]interface{}{
			"token":    token,
			"id":       u.ID,
			"nome":     u.Nome,
			"email":    u.Email,
			"username": u.Username,
			"foto_url": u.FotoURL,
			"cpf":      u.Cpf,
			"data_nascimento": u.DataNascimento,
		}, http.StatusOK)
	})

	mux.HandleFunc("POST /redefinir-senha", func(w http.ResponseWriter, r *http.Request) {
		var req dto.RedefinirSenhaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		if err := validation.ValidateStruct(req); err != nil {
			JSONError(w, err.Error(), http.StatusBadRequest)
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
		id, err := strconv.Atoi(idStr)
		if err != nil {
			JSONError(w, "ID inválido", http.StatusBadRequest)
			return
		}

		var req dto.UsuarioRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// Search current user data
		u, err := repo.BuscarPorID(r.Context(), id)
		if err != nil {
			JSONError(w, "Usuário não encontrado", http.StatusNotFound)
			return
		}

		// Update fields if provided
		if req.Nome != "" {
			u.Nome = req.Nome
		}
		if req.Email != "" {
			u.Email = req.Email
		}
		if req.FotoURL != "" {
			u.FotoURL = req.FotoURL
		}
		if req.Cpf != "" {
			u.Cpf = req.Cpf
		}
		if req.Username != "" {
			u.Username = req.Username
		}
		if req.DataNascimento != "" {
			t, err := time.Parse("2006-01-02", req.DataNascimento)
			if err == nil {
				u.DataNascimento = &t
			}
		}

		// Password Change Logic
		if req.NovaSenha != "" {
			if req.SenhaAtual == "" {
				JSONError(w, "Senha atual é obrigatória para alteração", http.StatusBadRequest)
				return
			}
			
			// Verify current password
			if !hash.VerificarHash(req.SenhaAtual, u.Senha) {
				JSONError(w, "Senha atual incorreta", http.StatusUnauthorized)
				return
			}

			// Generate new hash
			hashed, err := hash.GerarHash(req.NovaSenha)
			if err != nil {
				JSONError(w, "Erro ao processar nova senha", http.StatusInternalServerError)
				return
			}
			u.Senha = hashed
		}

		err = atualizarUC.Execute(r.Context(), *u)
		if err != nil {
			logger.Error("Erro ao atualizar usuário", zap.Int("id", id), zap.Error(err))
			JSONError(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		JSONSuccess(w, nil, http.StatusOK)
	})

	mux.HandleFunc("PUT /usuarios/{id}/meta", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			JSONError(w, "ID inválido", http.StatusBadRequest)
			return
		}

		var req dto.AtualizarMetaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			JSONError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		err = atualizarMetaUC.Execute(r.Context(), id, req.MetaPaginasSemana)
		if err != nil {
			logger.Error("Erro ao atualizar meta do usuário", zap.Int("id", id), zap.Error(err))
			JSONError(w, "Erro ao atualizar meta: "+err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, nil, http.StatusOK)
	})

	mux.HandleFunc("DELETE /usuarios/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			JSONError(w, "ID inválido", http.StatusBadRequest)
			return
		}
		if err := repo.Deletar(r.Context(), id); err != nil {
			JSONError(w, "Erro ao deletar usuário: "+err.Error(), http.StatusInternalServerError)
			return
		}
		JSONSuccess(w, nil, http.StatusOK)
	})
}
