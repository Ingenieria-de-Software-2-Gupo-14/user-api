package controller

import (
	"github.com/gin-gonic/gin"
	"ing-soft-2-tp1/cmd/api/database"
	. "ing-soft-2-tp1/cmd/api/models"
	"ing-soft-2-tp1/cmd/api/services"
	"net/http"
	"strconv"
)

type ControllerInterface interface {
	UserPost(context *gin.Context)
	UsersGet(context *gin.Context)
	UserGetById(context *gin.Context)
	UserDeleteById(context *gin.Context)
	dataBaseLength()
	addUser(user User)
	getUsers()
	getUser(id int)
	removeUser(id int)
}

// Controller struct that contains a database with users
type Controller struct {
	db database.Database[User]
}

// CreateController creates a controller
func CreateController() (controller Controller) {
	controller.db = database.CreateDatabase[User]()
	return controller
}

// dataBaseLength returns the amount of elements in the controllers database
func (controller Controller) dataBaseLength() (length int) {
	return controller.db.GetLen()
}

// addUser adds a user to the controller's database
func (controller Controller) addUser(user User) {
	controller.db.AddUser(user)
}

// getUsers return all elements of the database in List
func (controller Controller) getUsers() (users []User) {
	return controller.db.GetAllUsers()
}

// UsersPost adds the context's body as user in the database and sends a response to the api context. sends a 201 status code if its succesful and 400 if the title or description are missing
func (controller Controller) UsersPost(context *gin.Context) {
	var createUserRequest CreateUserRequest

	if err := context.Bind(&createUserRequest); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(StatusBadRequest, context.Request.URL.Path))
		return
	} else if createUserRequest.Title == "" {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(StatusMissingTitle, context.Request.URL.Path))
		return
	} else if createUserRequest.Description == "" {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(StatusMissingDescription, context.Request.URL.Path))
		return
	}
	user := services.CreateUser(controller.dataBaseLength(), createUserRequest.Title, createUserRequest.Description)
	controller.addUser(user)
	response := ResponseUser{User: user}
	context.JSON(201, response)
}

// UsersGet sends all users to the API context, even if there are none
func (controller Controller) UsersGet(context *gin.Context) {
	response := ResponseUsers{Users: controller.getUsers()}
	context.JSON(201, response)
}

// getUser return user from the database and "ok" bool value true in case that user was in database false in case it wasn't
func (controller Controller) getUser(id int) (user User, ok bool) {
	user, ok = controller.db.GetUser(id)
	return user, ok
}

// UserGetById sends response with the corresponding user with a status code 200, if the user isn't in the database it'll send a status code 404 not found
func (controller Controller) UserGetById(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(StatusInternalServerError, context.Request.URL.Path))
		return
	}
	user, ok := controller.getUser(id)
	if ok == false {
		context.JSON(http.StatusNotFound, services.CreateErrorResponse(StatusUserNotFound, context.Request.URL.Path))
		return
	}
	response := ResponseUser{User: user}
	context.JSON(http.StatusOK, response)
}

// removeUser removes user from database
func (controller Controller) removeUser(id int) {
	controller.db.DeleteUser(id)
}

// UserDeleteById removes user from database corresponding to id receive in context body, responds with code 204 "no content" in case of successful and 404 in case of user not found
func (controller Controller) UserDeleteById(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(StatusInternalServerError, context.Request.URL.Path))
		return
	}
	_, ok := controller.getUser(id)
	if ok == false {
		context.JSON(http.StatusNotFound, services.CreateErrorResponse(StatusUserNotFound, context.Request.URL.Path))
		return
	}
	controller.removeUser(id)
	context.JSON(http.StatusNoContent, nil)
}
