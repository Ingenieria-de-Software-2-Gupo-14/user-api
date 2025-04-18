package controller

import (
	. "ing-soft-2-tp1/internal/models"
	services "ing-soft-2-tp1/internal/services"
	"ing-soft-2-tp1/internal/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	DeleteUser(id int) error
	CreateUser(email string, password string, admin bool) (*User, error)
	GetUserById(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetAllUsers() (users []User, err error)
	ModifyUser(user *User) error
}

// UserController struct that contains a database with users
type UserController struct {
	service UserService
}

// CreateController creates a controller
func CreateController(service UserService) UserController {
	return UserController{service: service}
}

func (c UserController) RegisterUser(context *gin.Context) {
	var request CreateUserRequest

	if err := context.Bind(&request); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}

	if _, err := c.service.GetUserByEmail(request.Email); err == nil {
		context.JSON(http.StatusConflict, services.CreateErrorResponse(http.StatusConflict, context.Request.URL.Path))
		return
	}

	user, err := c.service.CreateUser(request.Email, request.Password, false)
	if err != nil {
		log.Println("Error creating user: ", err)
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}

	context.JSON(http.StatusCreated, ResponseUser{User: *user})
}

// UsersGet sends all users to the API context, even if there are none
func (c UserController) UsersGet(context *gin.Context) {
	users, err := c.service.GetAllUsers()
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

	user, ok := c.service.GetUserById(id)
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

	if err := controller.service.DeleteUser(id); err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	context.JSON(http.StatusNoContent, nil)
}

func (controller UserController) AdminsPost(context *gin.Context) {
	var createUserRequest CreateUserRequest

	if err := context.Bind(&createUserRequest); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}

	if _, err := controller.service.GetUserByEmail(createUserRequest.Email); err == nil {
		context.JSON(http.StatusConflict, services.CreateErrorResponse(http.StatusConflict, context.Request.URL.Path))
		return
	}

	user, err := controller.service.CreateUser(createUserRequest.Email, createUserRequest.Password, true)
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

	user, err := controller.service.GetUserByEmail(loginRequest.Email)
	if err != nil {
		log.Println("Error getting user by email: ", err)
		context.JSON(http.StatusUnauthorized, services.CreateErrorResponse(http.StatusUnauthorized, context.Request.URL.Path))
		return
	}

	log.Println("User found: ", user)
	if err := utils.CompareHashPassword(user.Password, loginRequest.Password); err != nil {
		log.Println("Error comparing password: ", err)
		context.JSON(http.StatusUnauthorized, services.CreateErrorResponse(http.StatusUnauthorized, context.Request.URL.Path))
		return
	}

	// Generar token JWT para la sesión
	token, err := utils.GenerateToken(user.Id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}

func (c UserController) ModifyUser(context *gin.Context) {
	var user User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}

	if err := c.service.ModifyUser(&user); err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}

	context.JSON(http.StatusOK, ResponseUser{User: user})
}

func (controller UserController) Health(context *gin.Context) {
	context.JSON(http.StatusOK, nil)
}
