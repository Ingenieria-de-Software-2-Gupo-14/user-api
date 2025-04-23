package router

import (
	"ing-soft-2-tp1/internal/controller"
	"ing-soft-2-tp1/internal/middleware"
	"ing-soft-2-tp1/internal/repositories"
	"ing-soft-2-tp1/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CreateRouter creates and return a Router with its corresponding end points
func CreateRouter(db *repositories.Database) *gin.Engine {
	r := gin.Default()
	userService := services.NewUserService(db)
	cont := controller.CreateController(userService)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // frontend address here
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // if you need cookies or auth headers
	}))
	r.GET("/health", cont.Health)
	r.POST("/users", cont.RegisterUser)
	r.POST("/admins", cont.RegisterAdmin)
	r.GET("/users", cont.UsersGet)
	r.POST("/users/modify", cont.ModifyUser)
	r.POST("/login", cont.UserLogin)
	r.GET("/users/:id", middleware.AuthMiddleware(), cont.UserGetById)
	r.DELETE("/users/:id", middleware.AuthMiddleware(), cont.UserDeleteById)
	r.PUT("/users/block/:id", cont.BlockUserById)
	r.PUT("/users/:id/location", cont.ModifyUserLocation)
	r.PUT("/admins/unblock/:id", cont.UnblockUserById)
	return r
}

func SetEnviroment(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}
