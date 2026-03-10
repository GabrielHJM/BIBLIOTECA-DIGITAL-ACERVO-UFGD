package main

import (
	"biblioteca-digital-api/config"
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()
	db := config.InitDB(cfg)
	defer db.Close()

	email := "gabriel@biblioteca.com"
	novaSenha := "123456"

	// Gerar hash bcrypt (cost 14 como no código original)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(novaSenha), 14)
	if err != nil {
		log.Fatal(err)
	}
	hashedPassword := string(hashedBytes)

	query := "UPDATE usuarios SET senha = $1 WHERE LOWER(email) = LOWER($2)"
	_, err = db.ExecContext(context.Background(), query, hashedPassword, email)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Senha do usuário %s resetada com sucesso para: %s\n", email, novaSenha)
}
