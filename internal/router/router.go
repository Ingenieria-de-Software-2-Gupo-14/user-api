package router

import (
	"ing-soft-2-tp1/internal/controller"
	"ing-soft-2-tp1/internal/repositories"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CreateRouter creates and return a Router with its corresponding end points
func CreateRouter(db *repositories.Database) *gin.Engine {
	r := gin.Default()
	cont := controller.CreateController(db)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"}, // frontend address here
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // if you need cookies or auth headers
	}))
	r.GET("/health", cont.Health)
	r.POST("/users", cont.UsersPost)
	r.POST("/admins", cont.AdminsPost)
	r.GET("/users", cont.UsersGet)
	r.POST("/users/modify", cont.ModifyUser)
	r.POST("/login", cont.UserLogin)
	r.GET("/users/:id", cont.UserGetById)
	r.DELETE("/users/:id", cont.UserDeleteById)
	return r
}

func SetEnviroment(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}
