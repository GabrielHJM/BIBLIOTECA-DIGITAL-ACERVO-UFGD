package config

import (
	"os"
)

type Config struct {
	DBUrl string
	Port  string
}

func Load() *Config {
	return &Config{
		DBUrl: getEnv("DATABASE_URL", "postgres://postgres@localhost:5432/postgres?sslmode=disable"),
		Port:  getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
