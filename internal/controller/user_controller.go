package controller

import (
	"context"
	. "ing-soft-2-tp1/internal/models"
	services "ing-soft-2-tp1/internal/services"
	"ing-soft-2-tp1/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	DeleteUser(ctx context.Context, id int) error
	CreateUser(ctx context.Context, request CreateUserRequest, admin bool) (*User, error)
	GetUserById(ctx context.Context, id int) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetAllUsers(ctx context.Context) (users []User, err error)
	ModifyUser(ctx context.Context, user *User) error
	BlockUser(ctx context.Context, id int) error
	ModifyLocation(ctx context.Context, id int, newLocation string) error
}

// UserController struct that contains a database with users
type UserController struct {
	service UserService
}

// CreateController creates a controller
func CreateController(service UserService) UserController {
	return UserController{service: service}
}

func (controller UserController) Health(context *gin.Context) {
	context.JSON(http.StatusOK, nil)
}

func (c UserController) RegisterUser(context *gin.Context) {
	var request CreateUserRequest
	ctx := context.Request.Context()
	if err := context.Bind(&request); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}

	if _, err := c.service.GetUserByEmail(ctx, request.Email); err == nil {
		context.JSON(http.StatusConflict, services.CreateErrorResponse(http.StatusConflict, context.Request.URL.Path))
		return
	}

	user, err := c.service.CreateUser(context.Request.Context(), request, false)
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}

	context.JSON(http.StatusCreated, ResponseUser{User: *user})
}

// UsersGet sends all users to the API context, even if there are none
func (c UserController) UsersGet(context *gin.Context) {
	users, err := c.service.GetAllUsers(context.Request.Context())
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	response := ResponseUsers{Users: users}
	context.JSON(http.StatusOK, response)
}

// UserGetById sends response with the corresponding user with a status code 200, if the user isn't in the database it'll send a status code 404 not found
func (c UserController) UserGetById(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}

	user, ok := c.service.GetUserById(context.Request.Context(), id)
	if ok != nil {
		context.JSON(http.StatusNotFound, services.CreateErrorResponse(StatusUserNotFound, context.Request.URL.Path))
		return
	}

	context.JSON(http.StatusOK, ResponseUser{User: *user})
}

// UserDeleteById removes user from database corresponding to id receive in context body, responds with code 204 "no content" in case of successful and 404 in case of user not found
func (controller UserController) UserDeleteById(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}

	if err := controller.service.DeleteUser(context.Request.Context(), id); err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	context.JSON(http.StatusNoContent, nil)
}

func (controller UserController) RegisterAdmin(context *gin.Context) {
	var createUserRequest CreateUserRequest

	if err := context.Bind(&createUserRequest); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}

	if _, err := controller.service.GetUserByEmail(context.Request.Context(), createUserRequest.Email); err == nil {
		context.JSON(http.StatusConflict, services.CreateErrorResponse(http.StatusConflict, context.Request.URL.Path))
		return
	}

	user, err := controller.service.CreateUser(context.Request.Context(), createUserRequest, true)
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}

	context.JSON(http.StatusCreated, ResponseUser{User: *user})
}

func (controller UserController) UserLogin(context *gin.Context) {
	var loginRequest LoginRequest
	if err := context.ShouldBindJSON(&loginRequest); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}

	user, err := controller.service.GetUserByEmail(context.Request.Context(), loginRequest.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, services.CreateErrorResponse(http.StatusNotFound, context.Request.URL.Path))
		return
	}

	if err := utils.CompareHashPassword(user.Password, loginRequest.Password); err != nil {
		context.JSON(http.StatusUnauthorized, services.CreateErrorResponse(http.StatusUnauthorized, context.Request.URL.Path))
		return
	}

	if user.BlockedUser == true {
		context.JSON(http.StatusForbidden, services.CreateErrorResponse(http.StatusForbidden, context.Request.URL.Path))
		return
	}

	token, err := GenerateToken(user.Id, user.Admin)
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}

	//Set cookie
	context.SetSameSite(http.SameSiteLaxMode)
	context.SetCookie("Authorization", token, 3600, "/", "", false, true)
	context.JSON(http.StatusOK, gin.H{"token": token})
}

func (c UserController) ModifyUser(context *gin.Context) {
	var user User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}

	if err := c.service.ModifyUser(context.Request.Context(), &user); err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}

	context.JSON(http.StatusOK, ResponseUser{User: user})
}

func (c UserController) ModifyUserLocation(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	var user User
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}
	if c.service.ModifyLocation(context, id, user.Location) != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	context.JSON(http.StatusOK, nil)
}

func (c UserController) BlockUserById(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	if c.service.BlockUser(context, id) != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	context.JSON(http.StatusOK, nil)
}

func (ct UserController) ValidateToken(c *gin.Context) {
	c.JSON(http.StatusOK, c.Request.Context().Value("jwt_info"))
}
