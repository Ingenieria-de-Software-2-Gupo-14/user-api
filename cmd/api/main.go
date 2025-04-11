package main

import (
	config2 "ing-soft-2-tp1/cmd/api/config"
	"ing-soft-2-tp1/cmd/api/router"
)

func main() {
	config := config2.LoadConfig() // lee las variables de entorno
	router.SetEnviroment(config.Environment)
	r := router.CreateRouter()
	if err := r.Run(config.Host + ":" + config.Port); err != nil {
		panic(err.Error())
	}
}
