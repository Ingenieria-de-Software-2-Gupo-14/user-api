package main

import (
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/auth"
	apiconfig "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/config"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/router"

	_ "github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

func main() {
	config := apiconfig.LoadConfig() // lee las variables de entorno

	auth.NewAuth(config)

	router.SetEnviroment(config.Environment)

	r := router.CreateRouter(config)
	if err := r.Run(":" + config.Port); err != nil {
		panic(err.Error())
	}
}
