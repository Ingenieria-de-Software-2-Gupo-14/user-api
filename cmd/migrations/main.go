package main

import (
	"ing-soft-2-tp1/internal/config"
	"log"
	"log/slog"

	"github.com/pressly/goose"
)

func main() {
	cfg := config.LoadConfig() // lee las variables de entorno

	db, err := config.CreateDatabase(cfg)
	if err != nil {
		log.Fatal("Error creating database", err)
	}

	// Run migrations
	if err := goose.Up(db, "internal/migrations"); err != nil {
		slog.Error("Error running migrations", slog.String("error", err.Error()))
	}
}
