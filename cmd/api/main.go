package main

import (
	apiconfig "ing-soft-2-tp1/internal/config"
	"ing-soft-2-tp1/internal/router"
)

func main() {
	config := apiconfig.LoadConfig() // lee las variables de entorno
	router.SetEnviroment(config.Environment)
	r := router.CreateRouter()
	if err := r.Run(config.Host + ":" + config.Port); err != nil {
		panic(err.Error())
	}
}
