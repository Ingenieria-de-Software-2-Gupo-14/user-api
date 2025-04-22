package main

import (
	apiconfig "ing-soft-2-tp1/internal/config"
	"ing-soft-2-tp1/internal/router"

	_ "github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

func main() {
	config := apiconfig.LoadConfig() // lee las variables de entorno
	userDatabase := apiconfig.SetupPostgresConnection()
	router.SetEnviroment(config.Environment)
	r := router.CreateRouter(userDatabase)
	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic(err.Error())
	}
}
