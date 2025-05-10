package router

import (
	"fmt"
	"net/http"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/config"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CreateRouter creates and return a Router with its corresponding end points
func CreateRouter(config config.Config) (*gin.Engine, error) {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // frontend address here
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // if you need cookies or auth headers
	}))

	deps, err := NewDependencies(&config)
	if err != nil {
		return nil, fmt.Errorf("error creating dependencies: %w", err)
	}

	r.GET("/health", func(ctx *gin.Context) {
		if err := deps.DB.Ping(); err != nil {
			utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"stats": deps.DB.Stats()})
	})

	// Auth routes
	auth := r.Group("/auth")
	auth.GET("/:provider", deps.Controllers.AuthController.BeginAuth)
	auth.GET("/:provider/callback", deps.Controllers.AuthController.CompleteAuth)
	auth.POST("/users", deps.Controllers.AuthController.Register(false))
	auth.POST("/users/verify", deps.Controllers.AuthController.VerifyRegistration)
	auth.POST("/admins", deps.Controllers.AuthController.Register(true))
	auth.POST("/login", deps.Controllers.AuthController.Login)
	auth.GET("/logout", deps.Controllers.AuthController.Logout)
	auth.PUT("/users/verify/resend", deps.Controllers.AuthController.ResendPin)

	// User routes
	r.GET("/users", deps.Controllers.UserController.UsersGet)
	r.POST("/users/modify", deps.Controllers.UserController.ModifyUser)
	r.GET("/users/:id", deps.Controllers.UserController.UserGetById)
	r.DELETE("/users/:id", deps.Controllers.UserController.UserDeleteById)
	r.PUT("/users/block/:id", deps.Controllers.UserController.BlockUserById)
	r.PUT("/users/:id/location", deps.Controllers.UserController.ModifyUserLocation)
	r.PUT("/users/:id/password", deps.Controllers.UserController.ModifyUserPasssword)
	return r, nil
}

func SetEnviroment(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}
