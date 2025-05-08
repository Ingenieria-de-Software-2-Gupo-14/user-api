package router

import (
	"database/sql"
	"log"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/auth"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/config"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/controller"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
)

type Dependencies struct {
	DB           *sql.DB
	Controllers  Controllers
	Services     Services
	Repositories Repositories
}

type Controllers struct {
	AuthController *auth.AuthController
	UserController *controller.UserController
}

type Services struct {
	UserService  services.UserService
	LoginService services.LoginAttemptService
	BlockService services.BlockedUserService
}

type Repositories struct {
	UserRepository  repositories.UserRepository
	LoginRepository repositories.LoginAttemptRepository
	BlockRepository repositories.BlockedUserRepository
}

func NewDependencies(cfg *config.Config) (*Dependencies, error) {
	db, err := cfg.CreateDatabase()
	if err != nil {
		log.Fatal("Error creating database", err)
	}

	// Repositories
	userRepo := repositories.CreateUserRepo(db)
	loginRepo := repositories.NewLoginAttemptRepository(db)
	blockRepo := repositories.NewBlockedUserRepository(db)
	// Services
	userService := services.NewUserService(userRepo)
	loginService := services.NewLoginAttemptService(loginRepo, blockRepo)
	blockService := services.NewBlockedUserService(blockRepo)

	// Controllers
	authController := auth.NewAuthController(userService, loginService)
	userController := controller.CreateController(userService)

	return &Dependencies{
		DB: db,
		Controllers: Controllers{
			AuthController: authController,
			UserController: userController,
		},
		Services: Services{
			UserService:  userService,
			LoginService: loginService,
			BlockService: blockService,
		},
		Repositories: Repositories{
			UserRepository:  userRepo,
			LoginRepository: loginRepo,
			BlockRepository: blockRepo,
		},
	}, nil

}
