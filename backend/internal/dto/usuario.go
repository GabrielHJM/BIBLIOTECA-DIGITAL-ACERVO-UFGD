package dto

type UsuarioRequest struct {
	Nome    string `json:"nome" validate:"required,min=3"`
	Email   string `json:"email" validate:"required,email"`
	Senha   string `json:"senha" validate:"required,min=6"`
	Tipo    int    `json:"tipo" validate:"required,oneof=1 2 3"`
	FotoURL string `json:"foto_url"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Senha string `json:"senha" validate:"required"`
}

type RedefinirSenhaRequest struct {
	Email string `json:"email" validate:"required,email"`
	Senha string `json:"senha" validate:"required,min=6"`
}

type AtualizarMetaRequest struct {
	MetaPaginasSemana int `json:"meta_paginas_semana"`
}
