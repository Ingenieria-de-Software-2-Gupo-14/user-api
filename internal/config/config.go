package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	Host        string
	Port        string
	Environment string

	// Database configuration
	DatabaseURL string

	// Secrets
	GoogleKey    string
	GoogleSecret string
	// JWT
	JWTSecret string
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// LoadConfig loads environment variables a Config Struct containing relevant variables
func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, loading environment variables from the system")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_QUERY"),
		)
	}

	return Config{
		Host:         getEnvOrDefault("HOST", "localhost"),
		Port:         getEnvOrDefault("PORT", "8080"),
		Environment:  getEnvOrDefault("ENVIRONMENT", "development"),
		DatabaseURL:  dbUrl,
		GoogleKey:    os.Getenv("GOOGLE_KEY"),
		GoogleSecret: os.Getenv("GOOGLE_SECRET"),
	}
}

func (config *Config) CreateDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
