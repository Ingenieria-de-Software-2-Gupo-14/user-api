package router

import (
	"fmt"
	"github.com/Ingenieria-de-Software-2-Gupo-14/go-core/pkg/telemetry"
	"net/http"
	"strings"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/config"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/middleware"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CreateRouter creates and return a Router with its corresponding end points
func CreateRouter(config config.Config) (*gin.Engine, error) {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return strings.HasSuffix(origin, ".vercel.app") || origin == "http://localhost:8081" || origin == "https://backoffice-seven-fawn.vercel.app"
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // if you need cookies or auth headers
	}))

	deps, err := NewDependencies(&config)
	if err != nil {
		return nil, fmt.Errorf("error creating dependencies: %w", err)
	}

	r.Use(telemetry.MetricsMiddleware(deps.Clients.TelemetryClient))

	r.GET("/health", func(ctx *gin.Context) {
		if err := deps.DB.Ping(); err != nil {
			utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"stats": deps.DB.Stats()})
	})

	//preflight options route
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(200)
	})

	// Auth routes
	auth := r.Group("/auth")
	auth.GET("/:provider", deps.Controllers.AuthController.BeginAuth)
	auth.GET("/:provider/callback", deps.Controllers.AuthController.CompleteAuth)
	auth.POST("/users", deps.Controllers.AuthController.Register)
	auth.POST("/users/verify", deps.Controllers.AuthController.VerifyRegistration)
	auth.POST("/admins", deps.Controllers.AuthController.RegisterAdmin)
	auth.POST("/login", deps.Controllers.AuthController.Login)
	auth.GET("/logout", deps.Controllers.AuthController.Logout)
	auth.PUT("/users/verify/resend", deps.Controllers.AuthController.ResendPin)

	// User routes
	r.GET("/users", middleware.AuthMiddleware(deps.Services.UserService), deps.Controllers.UserController.UsersGet)
	r.PUT("/users/:id", middleware.UserOrAdminMiddleware(deps.Services.UserService), deps.Controllers.UserController.ModifyUser)
	r.GET("/users/:id", middleware.AuthMiddleware(deps.Services.UserService), deps.Controllers.UserController.UserGetById)
	r.GET("/users/:id/notifications", deps.Controllers.UserController.GetUserNotifications)
	r.POST("/users/:id/notifications", deps.Controllers.UserController.SetUserNotifications)
	r.DELETE("/users/:id", deps.Controllers.UserController.UserDeleteById)
	r.PUT("/users/:id/block", middleware.AdminOnlyMiddleware(deps.Services.UserService), deps.Controllers.UserController.BlockUserById)
	r.PUT("/users/password", deps.Controllers.UserController.ModifyUserPasssword)
	r.POST("/users/reset/password", deps.Controllers.UserController.PasswordReset)
	r.GET("/users/reset/password", deps.Controllers.UserController.PasswordResetRedirect)
	r.POST("/users/notify", deps.Controllers.UserController.NotifyUsers)
	r.PUT("/users/:id/notifications/preference", deps.Controllers.UserController.ModifyNotifPreference)
	r.GET("/users/:id/notifications/preference", deps.Controllers.UserController.GetNotifPreferences)

	// Rules routes
	r.POST("/rules", deps.Controllers.UserController.AddRule)
	r.DELETE("/rules/:id", deps.Controllers.UserController.DeleteRule)
	r.GET("/rules", deps.Controllers.UserController.GetRules)
	r.PUT("/rules/:id", deps.Controllers.UserController.ModifyRule)
	r.GET("/rules/audit", deps.Controllers.UserController.GetAudits)

	//Ai Chat routes
	r.POST("/chat", middleware.AuthMiddleware(deps.Services.UserService), deps.Controllers.ChatController.SendMessage)
	r.GET("/chat", middleware.AuthMiddleware(deps.Services.UserService), deps.Controllers.ChatController.GetMessages)
	r.PUT("/chat/:message_id/rate", middleware.AuthMiddleware(deps.Services.UserService), deps.Controllers.ChatController.RateMessage)
	r.PUT("/chat/:message_id/feedback", middleware.AuthMiddleware(deps.Services.UserService), deps.Controllers.ChatController.FeedbackMessage)
	return r, nil
}

func SetEnviroment(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}
