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
	"log"

	_ "github.com/lib/pq"
)

func InitDB(cfg *Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Configuração do pool de conexões para melhor desempenho
	db.SetMaxOpenConns(25)                     // Limite máximo de conexões abertas
	db.SetMaxIdleConns(5)                      // Conexões inativas mantidas no pool
	db.SetConnMaxLifetime(5 * 60 * 1000000000) // Tempo de vida máximo (5 min) - usando nanosegundos (time.Minute seria melhor se importado)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Conexão com o banco de dados estabelecida e pool configurado.")
	return db
}
