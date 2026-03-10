package validation

import (
	"errors"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// Erros comuns de validação
var (
	ErrInvalidEmail     = errors.New("formato de email inválido")
	ErrPasswordTooShort = errors.New("a senha deve ter pelo menos 6 caracteres")
	ErrNameTooShort     = errors.New("o nome deve ter pelo menos 3 caracteres")
)

func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})
	return validate
}

// ValidateStruct valida um struct usando tags
func ValidateStruct(s interface{}) error {
	return GetValidator().Struct(s)
}

// Funções legadas para compatibilidade (podem ser removidas após refatoração total)

func ValidateEmail(email string) error {
	err := GetValidator().Var(email, "required,email")
	if err != nil {
		return ErrInvalidEmail
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(strings.TrimSpace(password)) < 6 {
		return ErrPasswordTooShort
	}
	return nil
}

func ValidateName(name string) error {
	if len(strings.TrimSpace(name)) < 3 {
		return ErrNameTooShort
	}
	return nil
}
