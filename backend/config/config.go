package config

import (
	"os"
)

type Config struct {
	DBUrl         string
	Port          string
	RedisURL      string
	RedisPassword string
}

func Load() *Config {
	return &Config{
		DBUrl:         getEnv("DATABASE_URL", "postgres://postgres@localhost:5432/postgres?sslmode=disable"),
		Port:          getEnv("PORT", "8080"),
		RedisURL:      getEnv("REDIS_URL", ""),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
