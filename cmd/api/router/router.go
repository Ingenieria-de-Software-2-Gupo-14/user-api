package router

import (
	"github.com/gin-gonic/gin"
	"ing-soft-2-tp1/cmd/api/controller"
)

// CreateRouter creates and return a Router with its corresponding end points
func CreateRouter() *gin.Engine {
	r := gin.Default()
	cont := controller.CreateController()
	r.POST("/users", cont.UsersPost)
	r.POST("/admins", cont.AdminsPost)
	r.GET("/users", cont.UsersGet)
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
