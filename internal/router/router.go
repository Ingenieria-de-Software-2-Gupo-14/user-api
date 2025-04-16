package router

import (
	"github.com/gin-gonic/gin"
	"ing-soft-2-tp1/internal/controller"
	"ing-soft-2-tp1/internal/database"
)

// CreateRouter creates and return a Router with its corresponding end points
func CreateRouter(db *database.Database) *gin.Engine {
	r := gin.Default()
	cont := controller.CreateController(db)
	r.POST("/users", cont.UsersPost)
	r.POST("/admins", cont.AdminsPost)
	r.GET("/users", cont.UsersGet)
	r.GET("/login/:email/:password", cont.UserGetByLogin)
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
