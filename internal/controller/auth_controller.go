package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/utils"

	"github.com/gin-gonic/gin"
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

	if user, err := ac.userRepo.GetUserByEmail(ctx, request.Email); err == nil {
		if !user.Verified {
			err := ac.userRepo.DeleteUser(ctx, user.Id)
			if err != nil {
				utils.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
				return
			}
		} else {
			utils.ErrorResponse(c, http.StatusConflict, "Email already exists")
			return
		}
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
	println(verification.VerificationPin)
	if verification.VerificationPin != parts[1] {
		utils.ErrorResponseWithErr(c, http.StatusBadRequest, err)
		return
	}
	errVerif := ac.userRepo.VerifyUser(ctx, verification.UserId)
	if errVerif != nil {
		utils.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
		return
	}
	err = ac.verificationService.DeleteByUserId(ctx, verification.UserId)
	if err != nil {
		utils.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
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

// GoogleAuth godoc
//
// @Summary      Authenticate with Google
// @Description  Authenticate a user using Google OAuth2
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body models.AuthRequest true "Google OAuth2 Token"
// @Success      200  {object}   map[string]interface{}
// @Failure      400  {object}   utils.HTTPError
// @Failure      401  {object}   utils.HTTPError
// @Failure      500  {object}   utils.HTTPError
// @Router       /auth/google [post]
func (ac *AuthController) GoogleAuth(c *gin.Context) {

	var request models.AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponseWithErr(c, http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()

	userInfo, err := models.ValidateGoogleToken(ctx, request.Token)
	if err != nil {
		utils.ErrorResponseWithErr(c, http.StatusUnauthorized, err)
		return
	}

	existingUser, err := ac.userRepo.GetUserByEmail(ctx, userInfo.Email)
	if err != nil {
		// User not found, create a new one
		if errors.Is(err, repositories.ErrNotFound) {

			id, err := ac.userRepo.CreateUser(ctx, userInfo)
			if err != nil {
				utils.ErrorResponseWithErr(c, http.StatusInternalServerError, err)
				return
			}
			existingUser = &models.User{
				Id:       id,
				Name:     userInfo.Name,
				Surname:  userInfo.Surname,
				Email:    userInfo.Email,
				Role:     "student",
				Verified: true,
			}
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

// @Summary      Sends a new Verification
// @Description  Sends a new Verification Pin to email saved in verification cookie
// @Tags         Auth
// @Accept       json
// @Produce      plain
// @Param        email  query     string  true  "Email"
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

// VerifyToken godoc
// @Summary      Verify a JWT token
// @Description  Verifies the JWT token and returns a success message if valid
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string  "Token is valid"
// @Failure      401  {object}  utils.HTTPError "Invalid or expired token"
// @Router       /auth/verify [get]
func (ac *AuthController) VerifyToken(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Token is valid"})
}
