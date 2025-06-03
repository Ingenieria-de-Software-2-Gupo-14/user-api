package main

import (
	"context"
	"log/slog"

	"github.com/Ingenieria-de-Software-2-Gupo-14/go-core/pkg/log"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/config"

	"github.com/pressly/goose"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig() // lee las variables de entorno

	db, err := cfg.CreateDatabase()
	if err != nil {
		log.Fatal(ctx, "Error creating database", slog.String("error", err.Error()))
	}
	// Run migrations
	if err := goose.Up(db, "internal/migrations"); err != nil {
		log.Error(ctx, "Error running migrations", slog.String("error", err.Error()))
	}

}
