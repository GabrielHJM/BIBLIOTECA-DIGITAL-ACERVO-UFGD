// package config

// import "os"

// func GetEnv(key, fallback string) string {
//     value := os.Getenv(key)
//     if value == "" {
//         return fallback
//     }
//     return value
// }
package config

type Config struct {
	DBUrl string
	Port  string
}

func Load() *Config {
	return &Config{
		DBUrl: "postgres://user:password@localhost:5432/dbname?sslmode=disable",
		Port:  "8080",
	}
}
