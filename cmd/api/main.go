package main

import (
	"ing-soft-2-tp1/internal/auth"
	apiconfig "ing-soft-2-tp1/internal/config"
	"ing-soft-2-tp1/internal/router"
	"log"

	_ "github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

func main() {
	config := apiconfig.LoadConfig() // lee las variables de entorno

	db, err := apiconfig.CreateDatabase(config)
	if err != nil {
		log.Fatal("Error creating database", err)
	}

	auth.NewAuth(config)

	router.SetEnviroment(config.Environment)

	r := router.CreateRouter(db)
	if err := r.Run(":" + config.Port); err != nil {
		panic(err.Error())
	}
}
