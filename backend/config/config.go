package config

import (
	"os"
)

type Config struct {
	DBUrl         string
	Port          string
	GeminiAPIKey  string
	RedisURL      string
	RedisPassword string
}

func Load() *Config {
	return &Config{
		DBUrl:         getEnv("DATABASE_URL", "postgres://postgres@localhost:5432/postgres?sslmode=disable"),
		Port:          getEnv("PORT", "8080"),
		GeminiAPIKey:  getEnv("GEMINI_API_KEY", ""),
		RedisURL:      getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
