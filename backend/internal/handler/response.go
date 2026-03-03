package handler

import (
	"encoding/json"
	"net/http"
)

// Response é a estrutura padrão para todas as respostas da API
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSONSuccess envia uma resposta de sucesso padronizada
func JSONSuccess(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
	})
}

// JSONError envia uma resposta de erro padronizada
// Dica: Para erros 500, a mensagem enviada ao cliente deve ser genérica
func JSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	displayMessage := message
	if statusCode == http.StatusInternalServerError {
		displayMessage = "Ocorreu um erro interno no servidor. Tente novamente mais tarde."
	}

	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   displayMessage,
	})
}
