package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/go-core/pkg/log"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"golang.org/x/net/context"
)

type AuthController struct {
	userRepo             services.UserService
	loginAttemptsService services.LoginAttemptService
	verificationService  services.VerificationService
}

func NewAuthController(userRepo services.UserService, loginAttemptsService services.LoginAttemptService, verificationService services.VerificationService) *AuthController {
	return &AuthController{
		userRepo:             userRepo,
		loginAttemptsService: loginAttemptsService,
		verificationService:  verificationService,
	}
}

// @Summary      Register a new user
// @Description  Registers a new user
// @Tags         Auth
// @Accept       json
// @Produce      plain
// @Param        request body models.CreateUserRequest true "User Registration Details"
// @Success      201  {object}  map[string]int  "User created successfully"
// @Failure      400  {object}  utils.HTTPError "Invalid request format"
// @Failure      409  {object}  utils.HTTPError "Email already exists"
// @Failure      500  {object}  utils.HTTPError "Internal server error"
// @Router       /auth/users [post]
func (ac *AuthController) Register(c *gin.Context) {
	var request models.CreateUserRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	if _, err := ac.userRepo.GetUserByEmail(ctx, request.Email); err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Email already exists")
		return
	}

	id, err := ac.userRepo.CreateUser(ctx, request)
	if err != nil {
		utils.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
		return
	}

	if err := ac.verificationService.SendVerificationEmail(ctx, id, request.Email); err != nil {
		utils.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary      Register a new Admin
// @Description  Registers a new Admin
// @Tags         Auth
// @Accept       json
// @Produce      plain
// @Param        request body models.CreateUserRequest true "User Registration Details"
// @Success      201  {object}  map[string]int  "User created successfully"
// @Failure      400  {object}  utils.HTTPError "Invalid request format"
// @Failure      409  {object}  utils.HTTPError "Email already exists"
// @Failure      500  {object}  utils.HTTPError "Internal server error"
// @Router       /auth/admins [post]
func (ac *AuthController) RegisterAdmin(c *gin.Context) {
	var request models.CreateUserRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	if _, err := ac.userRepo.GetUserByEmail(ctx, request.Email); err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Email already exists")
		return
	}

	request.Role = "admin"
	request.Verified = true

	id, err := ac.userRepo.CreateUser(c.Request.Context(), request)
	if err != nil {
		utils.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary      Verify a new user's registration
// @Description  Verify the new user's registration using a Pin sent to the users email
// @Tags         Auth
// @Accept       json
// @Produce      plain
// @Param        request body models.EmailVerifiaction true "Email Verification Details"
// @Success      201  {object}  map[string]int  "User Verified and created successfully"
// @Failure      400  {object}  utils.HTTPError "Invalid request format"
// @Failure      500  {object}  utils.HTTPError "Internal server error"
// @Router       /auth/users/verify [post]
func (ac *AuthController) VerifyRegistration(c *gin.Context) {
	var request models.EmailVerifiaction
	if err := c.Bind(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format")
	}

	parts := strings.Split(request.VerificationPin, "-")
	if len(parts) != 2 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid verification pin format")
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid verification pin format")
		return
	}

	ctx := c.Request.Context()
	verification, err := ac.verificationService.GetVerification(ctx, id)
	if err != nil {
		utils.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
		return
	}

	if time.Now().After(verification.PinExpiration) {
		utils.ErrorResponseWithErr(c, http.StatusGone, err)
		return
	}

	if verification.VerificationPin != request.VerificationPin {
		utils.ErrorResponseWithErr(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User verified and created successfully"})
}

// finishAuth finish the authentication process by generating a token and setting it in the cookie
// and redirecting to the home page
func (ac *AuthController) finishAuth(ctx *gin.Context, user models.User) {
	if user.Blocked {
		utils.ErrorResponse(ctx, http.StatusForbidden, "User is blocked")
		return
	}

	token, err := models.GenerateToken(user.Id, user.Email, user.Name, user.Role)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}

	ac.loginAttemptsService.AddLoginAttempt(ctx, user.Id, ctx.Request.RemoteAddr, ctx.Request.UserAgent(), true)

	ctx.SetSameSite(http.SameSiteLaxMode)

	ctx.SetCookie("Authorization", token, 3600, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// Login godoc
// @Summary      User login
// @Description  Logs in a user with email and password, returning user information and a JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.LoginRequest  true  "Login Credentials"
// @Success      200      {object}  map[string]interface{}  "Successful login returns user and token"
// @Failure      400      {object}  utils.HTTPError   "Invalid request format or input validation failed"
// @Failure      401      {object}  utils.HTTPError   "Invalid email or password"
// @Failure      403      {object}  utils.HTTPError   "User is blocked"
// @Failure      500      {object}  utils.HTTPError   "Internal server error (e.g., token generation failed)"
// @Router       /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponseWithErr(c, http.StatusBadRequest, err)
		return
	}

	user, err := ac.userRepo.GetUserByEmail(c.Request.Context(), request.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err := utils.CompareHashPassword(user.Password, request.Password); err != nil {
		ac.loginAttemptsService.AddLoginAttempt(c, user.Id, c.Request.RemoteAddr, c.Request.UserAgent(), false)
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if user.Blocked {
		utils.ErrorResponse(c, http.StatusForbidden, "User is blocked")
		return
	}

	if !user.Verified {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User is not verified")
		return
	}

	ac.finishAuth(c, *user)
}

// BeginAuth godoc
//
// @Summary      Begin authentication
// @Description  Begin authentication with the specified provider
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        provider  path      string  true  "Provider name"
// @Success      200  {object}   map[string]interface{}
// @Failure      400  {object}   utils.HTTPError
// @Failure      401  {object}   utils.HTTPError
// @Failure      500  {object}   utils.HTTPError
// @Router       /auth/{provider} [get]
func (ac *AuthController) BeginAuth(c *gin.Context) {
	provider := c.Param("provider")
	ctx := context.WithValue(c.Request.Context(), "provider", provider)
	r := c.Request.WithContext(ctx)

	gothic.BeginAuthHandler(c.Writer, r)

}

// CompleteAuth godoc
//
// @Summary      Complete authentication
// @Description  Complete authentication with the specified provider
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        provider  path      string  true  "Provider name"
// @Success      200  {object}   map[string]interface{}
// @Failure      400  {object}   utils.HTTPError
// @Failure      401  {object}   utils.HTTPError
// @Failure      500  {object}   utils.HTTPError
// @Router       /auth/{provider}/callback [get]
func (ac *AuthController) CompleteAuth(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), "provider", c.Param("provider"))
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request.WithContext(ctx))
	if err != nil {
		utils.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
		return
	}

	existingUser, err := ac.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		// User not found, create a new one
		if errors.Is(err, repositories.ErrNotFound) {
			newUser := models.CreateUserRequest{
				Name:     user.Name,
				Surname:  user.LastName,
				Email:    user.Email,
				Role:     "student",
				Verified: true,
			}

			id, err := ac.userRepo.CreateUser(ctx, newUser)
			if err != nil {
				utils.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
				return
			}

			existingUser.Id = id
			existingUser.Name = user.Name
			existingUser.Surname = user.LastName
			existingUser.Email = user.Email
		} else {
			utils.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
			return
		}
	}

	ac.finishAuth(c, *existingUser)
}

// Logout godoc
//
// @Summary      Logout
// @Description  Logout the user by clearing the cookie
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      307  {string} string "Redirected to home page"
// @Router       /auth/logout [get]
func (ac *AuthController) Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func getAuthToken(c *gin.Context) string {
	auth, _ := c.Cookie("Authorization")
	if auth == "" {
		if parts := strings.Fields(c.GetHeader("Authorization")); len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			auth = parts[1]
		}
	}
	return auth
}

func (ac *AuthController) AuthMiddlewarefunc(ctx *gin.Context) {
	tokenStr := getAuthToken(ctx)
	claims, err := models.ParseToken(tokenStr)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusUnauthorized, err)
		ctx.Abort()
		return
	}

	uId, _ := strconv.Atoi(claims.Subject)
	blocked, err := ac.userRepo.IsUserBlocked(ctx.Request.Context(), uId)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		ctx.Abort()
		return
	}

	if blocked {
		utils.ErrorResponse(ctx, http.StatusForbidden, "User is blocked")
		// Make Cookie Expire
		ctx.SetCookie("Authorization", "", -1, "/", "", false, true)
		ctx.Abort()
		return
	}

	// Refresh token if it is about to expire
	if claims.ExpiresAt < time.Now().Add(time.Minute*5).Unix() {
		newToken, err := models.GenerateToken(uId, claims.Email, claims.Name, claims.Role)
		if err != nil {
			log.Error(ctx, "Error generating token", "error", err.Error())
		} else {
			ctx.SetCookie("Authorization", newToken, 3600, "/", "", false, true)
		}
	}

	ctx.Set("claims", claims)
	ctx.Next()
}

// @Summary      Sends a new Verification
// @Description  Sends a new Verification Pin to email saved in verification cookie
// @Tags         Auth
// @Accept       json
// @Produce      plain
// @Success      200  {object}  nil  "New Pin sent successfully"
// @Failure      500  {object}  utils.HTTPError "Internal server error"
// @Router       /auth/users/verify/resend [put]
func (ac *AuthController) ResendPin(c *gin.Context) {

	//get email from query params
	email := c.Query("email")
	if email == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Email is required")
		return
	}

	//Get user by email
	user, err := ac.userRepo.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid verification")
			return
		}
		utils.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
		return
	}

	err = ac.verificationService.UpdatePin(c.Request.Context(), user.Id, email)
	if err != nil {
		utils.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
