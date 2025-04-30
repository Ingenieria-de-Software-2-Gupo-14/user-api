package router

import (
	"database/sql"
	"ing-soft-2-tp1/internal/auth"
	"ing-soft-2-tp1/internal/controller"
	"ing-soft-2-tp1/internal/repositories"
	"ing-soft-2-tp1/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CreateRouter creates and return a Router with its corresponding end points
func CreateRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	userRepo := repositories.CreateUserRepo(db)
	userService := services.NewUserService(userRepo)
	cont := controller.CreateController(userService)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // frontend address here
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // if you need cookies or auth headers
	}))

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	auth.AddAuthRoutes(r, userRepo)

	r.POST("/users", cont.RegisterUser)
	r.POST("/admins", cont.RegisterAdmin)
	r.GET("/users", cont.UsersGet)
	r.POST("/users/modify", cont.ModifyUser)
	r.POST("/login", cont.UserLogin)
	r.GET("/users/:id", auth.AuthMiddleware(userRepo), cont.UserGetById)
	r.DELETE("/users/:id", auth.AuthMiddleware(userRepo), cont.UserDeleteById)
	r.PUT("/users/block/:id", cont.BlockUserById)
	r.PUT("/users/:id/location", cont.ModifyUserLocation)
	r.PUT("/users/:id/privacy", cont.ModifyUserPrivacy)
	r.GET("/users/:id/privacy", cont.UserGetPrivacy)
	return r
}

func SetEnviroment(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}
