package validation

import (
	"errors"
	"net/mail"
	"strings"
)

// Erros comuns de validação
var (
	ErrInvalidEmail     = errors.New("formato de email inválido")
	ErrPasswordTooShort = errors.New("a senha deve ter pelo menos 6 caracteres")
	ErrNameTooShort     = errors.New("o nome deve ter pelo menos 3 caracteres")
)

// ValidateEmail verifica se o email é válido
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}
	return nil
}

// ValidatePassword verifica a força da senha
func ValidatePassword(password string) error {
	if len(strings.TrimSpace(password)) < 6 {
		return ErrPasswordTooShort
	}
	return nil
}

// ValidateName verifica o nome
func ValidateName(name string) error {
	if len(strings.TrimSpace(name)) < 3 {
		return ErrNameTooShort
	}
	return nil
}
