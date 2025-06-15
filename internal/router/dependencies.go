package router

import (
	"database/sql"
	"os"

	"github.com/Ingenieria-de-Software-2-Gupo-14/go-core/pkg/telemetry"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/config"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/controller"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/sendgrid/sendgrid-go"
)

type Dependencies struct {
	DB           *sql.DB
	Controllers  Controllers
	Services     Services
	Repositories Repositories
	Clients      Clients
}

type Controllers struct {
	AuthController *controller.AuthController
	UserController *controller.UserController
	ChatController *controller.ChatController
}

type Services struct {
	UserService  services.UserService
	LoginService services.LoginAttemptService
}

type Repositories struct {
	UserRepository  repositories.UserRepository
	LoginRepository repositories.LoginAttemptRepository
	BlockRepository repositories.BlockedUserRepository
}

type Clients struct {
	TelemetryClient telemetry.Client
}

func NewDependencies(cfg *config.Config) (*Dependencies, error) {
	db, err := cfg.CreateDatabase()
	if err != nil {
		return nil, err
	}

	// Repositories
	userRepo := repositories.CreateUserRepo(db)
	loginRepo := repositories.NewLoginAttemptRepository(db)
	blockRepo := repositories.NewBlockedUserRepository(db)
	verificationRepo := repositories.CreateVerificationRepo(db)
	rulesRepo := repositories.CreateRulesRepo(db)
	chatRepo := repositories.CreateChatsRepo(db)
	// Services
	userService := services.NewUserService(userRepo, blockRepo)
	loginService := services.NewLoginAttemptService(loginRepo, blockRepo)
	verificationService := services.NewVerificationService(verificationRepo, sendgrid.NewSendClient(os.Getenv("EMAIL_API_KEY")))
	rulesService := services.NewRulesService(rulesRepo)
	chatService := services.NewChatsService(chatRepo)

	// Controllers
	authController := controller.NewAuthController(userService, loginService, verificationService)
	userController := controller.CreateController(userService, rulesService)
	chatController := controller.NewChatsController(chatService)

	// Clients
	telemetryClient, err := cfg.CreateDatadogClient()
	if err != nil {
		return nil, err
	}

	return &Dependencies{
		DB: db,
		Controllers: Controllers{
			AuthController: authController,
			UserController: userController,
			ChatController: chatController,
		},
		Services: Services{
			UserService:  userService,
			LoginService: loginService,
		},
		Repositories: Repositories{
			UserRepository:  userRepo,
			LoginRepository: loginRepo,
			BlockRepository: blockRepo,
		},
		Clients: Clients{
			TelemetryClient: telemetryClient,
		},
	}, nil

}
