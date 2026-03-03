package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Terminate all connections to postgres to allow using it as a template
	_, err = db.Exec(`SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = 'postgres' AND pid <> pg_backend_pid();`)
	if err != nil {
		fmt.Println("Terminate error:", err)
	}

	_, err = db.Exec("CREATE DATABASE biblioteca_digital WITH TEMPLATE postgres;")
	if err != nil {
		fmt.Println("Create error:", err)
	}
	fmt.Println("Done")
}
