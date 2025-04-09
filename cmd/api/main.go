package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Leer variables de entorno
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "local"
	}

	log.Printf("INFO: Starting server in %s mode on %s:%s", environment, host, port)

	// Configurar Gin
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	// Iniciar el servidor
	r.Run(host + ":" + port)
}
