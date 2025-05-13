package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Ingenieria-de-Software-2-Gupo-14/go-core/pkg/telemetry"
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

	// Telemetry configuration
	DatadogClientType string
	DatadogHost       string
	DatadogStatsdPort string

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
		Host:              getEnvOrDefault("HOST", "localhost"),
		Port:              getEnvOrDefault("PORT", "8080"),
		Environment:       getEnvOrDefault("ENVIRONMENT", "development"),
		DatabaseURL:       dbUrl,
		GoogleKey:         os.Getenv("GOOGLE_KEY"),
		GoogleSecret:      os.Getenv("GOOGLE_SECRET"),
		DatadogClientType: getEnvOrDefault("DD_CLIENT_TYPE", "default"),
		DatadogHost:       getEnvOrDefault("DD_HOST", "localhost"),
		DatadogStatsdPort: getEnvOrDefault("DD_STATSD_PORT", "8125"),
	}
}

func (config *Config) CreateDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Set timezone to UTC for all sessions
	_, err = db.Exec("SET timezone TO 'UTC';")
	if err != nil {
		return nil, fmt.Errorf("failed to set timezone to UTC: %w", err)
	}

	return db, nil
}

func (config *Config) CreateDatadogClient() (telemetry.Client, error) {
	switch config.DatadogClientType {
	case "api":
		return telemetry.NewDatadogAPI()

	case "statsd", "agent":
		return telemetry.NewDatadog(config.DatadogHost + ":" + config.DatadogStatsdPort)

	default:
		return nil, nil
	}
}
