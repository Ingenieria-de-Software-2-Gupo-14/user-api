package main

import (
	"context"
	"log/slog"

	"github.com/Ingenieria-de-Software-2-Gupo-14/go-core/pkg/log"
	_ "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/docs"
	apiconfig "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/config"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/router"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

// @title           User API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	config := apiconfig.LoadConfig() // lee las variables de entorno
	ctx := context.Background()
	goth.UseProviders(
		google.New(config.GoogleKey, config.GoogleSecret, "http://localhost:8080/auth/google/callback"),
	)

	router.SetEnviroment(config.Environment)

	r, err := router.CreateRouter(config)
	if err != nil {
		log.Fatal(ctx, "Error creating router", slog.String("error", err.Error()))
	}
	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := r.Run(":" + config.Port); err != nil {
		log.Fatal(ctx, "Error running router", slog.String("error", err.Error()))
	}
}
