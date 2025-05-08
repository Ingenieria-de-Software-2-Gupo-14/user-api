package main

import (
	"log"
	"log/slog"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/config"

	"github.com/pressly/goose"
)

func main() {
	cfg := config.LoadConfig() // lee las variables de entorno

	db, err := cfg.CreateDatabase()
	if err != nil {
		log.Fatal("Error creating database", err)
	}

	// Run migrations
	if err := goose.Up(db, "internal/migrations"); err != nil {
		slog.Error("Error running migrations", slog.String("error", err.Error()))
	}
}
