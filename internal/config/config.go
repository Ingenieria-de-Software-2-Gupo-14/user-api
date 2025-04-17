package config

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"ing-soft-2-tp1/internal/database"
	"log"
	"os"
)

type Config struct {
	Host        string
	Port        string
	Environment string
}

// LoadConfig loads environment variables a Config Struct containing relevant variables
func LoadConfig() Config {

	return Config{os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("ENVIROMENT")}
}

func SetupPostgresConnection() *database.Database {
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbUser := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName))
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Run migrations
	if err := goose.Up(db, "internal/migrations"); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	userDatabase := database.CreateDatabase(db)
	return userDatabase
}
