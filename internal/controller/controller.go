package controller

import (
	"github.com/gin-gonic/gin"
	. "ing-soft-2-tp1/internal/models"
	services "ing-soft-2-tp1/internal/services"
	"net/http"
	"strconv"
)

type Database interface {
	GetUser(id int) (*User, error)
	GetAllUsers() ([]User, error)
	DeleteUser(id int) error
	AddUser(user *User) (int, error)
	GetUserByEmailAndPassword(email string, password string) (*User, error)
	ContainsUserByEmail(email string) bool
	ModifyUser(user *User) error
}

// Controller struct that contains a database with users
type Controller struct {
	db Database
}

// CreateController creates a controller
func CreateController(db Database) (controller Controller) {
	controller.db = db
	return controller
}

// UsersPost adds the context's body as user in the database and sends a response to the api context. sends a 201 status code if its succesful and 400 if the title or description are missing
func (controller Controller) UsersPost(context *gin.Context) {
	var createUserRequest CreateUserRequest

	if err := context.Bind(&createUserRequest); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}
	if services.ContainsUserByEmail(controller.db, createUserRequest.Email) {
		context.JSON(http.StatusConflict, services.CreateErrorResponse(http.StatusConflict, context.Request.URL.Path))
	}
	user := services.CreateUser(0, createUserRequest.Email, createUserRequest.Password)
	_, err := services.AddUserToDatabase(controller.db, &user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	response := ResponseUser{User: user}
	context.JSON(201, response)
}

// UsersGet sends all users to the API context, even if there are none
func (controller Controller) UsersGet(context *gin.Context) {
	users, err := services.GetAllUsersFromDatabase(controller.db)
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	response := ResponseUsers{Users: users}
	context.JSON(201, response)
}

// UserGetById sends response with the corresponding user with a status code 200, if the user isn't in the database it'll send a status code 404 not found
func (controller Controller) UserGetById(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	user, ok := services.GetUserFromDatabase(controller.db, id)
	if ok != nil {
		context.JSON(http.StatusNotFound, services.CreateErrorResponse(StatusUserNotFound, context.Request.URL.Path))
		return
	}
	response := ResponseUser{User: *user}
	context.JSON(http.StatusOK, response)
}

// UserDeleteById removes user from database corresponding to id receive in context body, responds with code 204 "no content" in case of successful and 404 in case of user not found
func (controller Controller) UserDeleteById(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	services.RemoveUserFromDatabase(controller.db, id)
	context.JSON(http.StatusNoContent, nil)
}

func (controller Controller) AdminsPost(context *gin.Context) {
	var createUserRequest CreateUserRequest

	if err := context.Bind(&createUserRequest); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}
	if services.ContainsUserByEmail(controller.db, createUserRequest.Email) {
		context.JSON(http.StatusConflict, services.CreateErrorResponse(http.StatusConflict, context.Request.URL.Path))
	}
	user := services.CreateAdminUser(0, createUserRequest.Email, createUserRequest.Password)
	id, err := services.AddUserToDatabase(controller.db, &user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	println(id)
	user.Id = id
	response := ResponseUser{User: user}
	context.JSON(201, response)

}

func (controller Controller) UserLogin(context *gin.Context) {
	var loginRequest LoginRequest

	if err := context.ShouldBindJSON(&loginRequest); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}
	user, ok := services.GetUserFromDatabaseByEmailAndPassword(controller.db, loginRequest.Email, loginRequest.Password)
	if ok != nil {
		context.JSON(http.StatusNotFound, services.CreateErrorResponse(StatusUserNotFound, context.Request.URL.Path))
		return
	}
	response := ResponseUser{User: *user}
	context.JSON(http.StatusOK, response)
}

func (controller Controller) ModifyUser(context *gin.Context) {
	var user User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}

	ok := services.ModifyUser(controller.db, &user)
	if ok != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	context.JSON(http.StatusNoContent, nil)
}

func (controller Controller) Health(context *gin.Context) {
	context.JSON(http.StatusOK, nil)
}
