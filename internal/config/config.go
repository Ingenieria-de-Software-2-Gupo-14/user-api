package config

import (
	"database/sql"
	"ing-soft-2-tp1/internal/repositories"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
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

func SetupPostgresConnection() *repositories.Database {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL")) //fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName))
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Run migrations
	if err := goose.Up(db, "internal/migrations"); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	log.Println("Connected to database", db.Ping())
	userDatabase := repositories.CreateDatabase(db)
	return userDatabase
}
