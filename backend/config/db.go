// package config

// import (
//     "database/sql"
//     _ "github.com/lib/pq"
//     "log"
// )

// func ConnectDB() (*sql.DB, error) {
//     dsn := GetEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/db?sslmode=disable")
//     db, err := sql.Open("postgres", dsn)
//     if err != nil {
//         log.Fatal(err)
//     }
//     return db, nil
// }
package config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func InitDB(cfg *Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
