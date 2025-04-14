package controller

import (
	"github.com/gin-gonic/gin"
	. "ing-soft-2-tp1/internal/database"
	. "ing-soft-2-tp1/internal/models"
	services "ing-soft-2-tp1/internal/services"
	"net/http"
	"strconv"
)

// Controller struct that contains a database with users
type Controller struct {
	db Database[User]
}

// CreateController creates a controller
func CreateController() (controller Controller) {
	controller.db = CreateDatabase[User]()
	return controller
}

// dataBaseLength returns the amount of elements in the controllers database
func (controller Controller) dataBaseLength() (length int) {
	return controller.db.GetLen()
}

// UsersPost adds the context's body as user in the database and sends a response to the api context. sends a 201 status code if its succesful and 400 if the title or description are missing
func (controller Controller) UsersPost(context *gin.Context) {
	var createUserRequest CreateUserRequest

	if err := context.Bind(&createUserRequest); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}
	user := services.CreateUser(controller.dataBaseLength(), createUserRequest.Email, createUserRequest.Password)
	services.AddUserToDatabase(&controller.db, user)
	response := ResponseUser{User: user}
	context.JSON(201, response)
}

// UsersGet sends all users to the API context, even if there are none
func (controller Controller) UsersGet(context *gin.Context) {
	response := ResponseUsers{Users: services.GetAllUsersFromDatabase(&controller.db)}
	context.JSON(201, response)
}

// UserGetById sends response with the corresponding user with a status code 200, if the user isn't in the database it'll send a status code 404 not found
func (controller Controller) UserGetById(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	user, ok := services.GetUserFromDatabase(&controller.db, id)
	if ok == false {
		context.JSON(http.StatusNotFound, services.CreateErrorResponse(StatusUserNotFound, context.Request.URL.Path))
		return
	}
	response := ResponseUser{User: user}
	context.JSON(http.StatusOK, response)
}

// UserDeleteById removes user from database corresponding to id receive in context body, responds with code 204 "no content" in case of successful and 404 in case of user not found
func (controller Controller) UserDeleteById(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, services.CreateErrorResponse(http.StatusInternalServerError, context.Request.URL.Path))
		return
	}
	_, ok := services.GetUserFromDatabase(&controller.db, id)
	if ok == false {
		context.JSON(http.StatusNotFound, services.CreateErrorResponse(StatusUserNotFound, context.Request.URL.Path))
		return
	}
	services.RemoveUserFromDatabase(&controller.db, id)
	context.JSON(http.StatusNoContent, nil)
}

func (controller Controller) AdminsPost(context *gin.Context) {
	var createUserRequest CreateUserRequest

	if err := context.Bind(&createUserRequest); err != nil {
		context.JSON(http.StatusBadRequest, services.CreateErrorResponse(http.StatusBadRequest, context.Request.URL.Path))
		return
	}
	user := services.CreateAdminUser(controller.dataBaseLength(), createUserRequest.Email, createUserRequest.Password)
	services.AddUserToDatabase(&controller.db, user)
	response := ResponseUser{User: user}
	context.JSON(201, response)

}
