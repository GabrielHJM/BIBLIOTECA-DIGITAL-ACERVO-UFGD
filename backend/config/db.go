// package config

// import (
//     "database/sql"
//     _ "github.com/lib/pq"
//     "log"
// )

//	func ConnectDB() (*sql.DB, error) {
//	    dsn := GetEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/db?sslmode=disable")
//	    db, err := sql.Open("postgres", dsn)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    return db, nil
//	}
package config

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func InitDB(cfg *Config) *sql.DB {
	// Log sanitized DB info
	u, err := url.Parse(cfg.DBUrl)
	if err == nil {
		log.Printf("Tentando conectar ao banco de dados em: %s:%s%s", u.Hostname(), u.Port(), u.Path)
	} else {
		log.Printf("URL do banco de dados inválida ou vazia.")
	}

	// 1. Tentar conectar ao banco principal
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Printf("Erro ao abrir banco de dados: %v", err)
	}

	// Tentar um Ping para ver se o banco existe
	err = db.Ping()
	if err != nil && (strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "database") ) {
		log.Printf("Banco de dados não encontrado. Tentando criar automaticamente...")

		// 2. Tentar criar o banco se ele não existir
		if err := createDatabase(cfg.DBUrl); err != nil {
			log.Printf("Aviso: Não foi possível criar o banco automaticamente: %v", err)
			log.Println("Por favor, crie o banco 'BibliotecaDigital_BD' manualmente no pgAdmin.")
		} else {
			log.Println("Banco de dados criado com sucesso.")
			// Reabrir conexão após criação
			db, err = sql.Open("postgres", cfg.DBUrl)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else if err != nil {
		log.Printf("Erro de conexão detectado: %v", err)
	}

	// Configuração do pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Validação final
	if err := db.Ping(); err != nil {
		log.Printf("Falha crítica ao conectar ao banco de dados: %v", err)
		log.Fatal("Certifique-se de que o PostgreSQL está rodando e as credenciais no .env estão corretas.")
	}

	log.Println("Conexão com o banco de dados estabelecida.")
	return db
}

func createDatabase(dbUrl string) error {
	u, err := url.Parse(dbUrl)
	if err != nil {
		return err
	}

	// Nome do banco que queremos criar
	dbName := strings.TrimPrefix(u.Path, "/")

	// Alterar a URL para conectar ao banco padrão 'postgres'
	u.Path = "/postgres"
	postgresUrl := u.String()

	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		return err
	}
	defer db.Close()

	// Executar comando de criação
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", dbName))
	return err
}

