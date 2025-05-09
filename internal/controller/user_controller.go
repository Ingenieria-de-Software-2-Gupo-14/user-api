package controller

import (
	"net/http"
	"strconv"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	services "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/utils"

	"github.com/gin-gonic/gin"
)

// UserController struct that contains a database with users
type UserController struct {
	service services.UserService
}

// CreateController creates a controller
func CreateController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// UsersGet godoc
// @Summary      Get all users
// @Description  Returns a list of all users in the system
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string][]models.User  "List of users"
// @Failure      500  {object}  utils.HTTPError          "Internal server error"
// @Router       /users [get]
func (c UserController) UsersGet(context *gin.Context) {
	users, err := c.service.GetAllUsers(context.Request.Context())
	if err != nil {
		utils.ErrorResponseWithErr(context, http.StatusInternalServerError, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": users})
}

// UserGetById godoc
// @Summary      Get user by ID
// @Description  Returns a specific user by their ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]models.User  "User data"
// @Failure      400  {object}  utils.HTTPError        "Invalid user ID format"
// @Failure      404  {object}  utils.HTTPError        "User not found"
// @Router       /users/{id} [get]
func (c UserController) UserGetById(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	user, err := c.service.GetUserById(context.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(context, http.StatusNotFound, "User not found")
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": user})
}

// UserDeleteById godoc
// @Summary      Delete user by ID
// @Description  Removes a user from the database
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      204  {object}  nil  "User successfully deleted"
// @Failure      400  {object}  utils.HTTPError  "Invalid user ID format"
// @Failure      500  {object}  utils.HTTPError  "Internal server error"
// @Router       /users/{id} [delete]
func (controller UserController) UserDeleteById(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	if err := controller.service.DeleteUser(context.Request.Context(), id); err != nil {
		utils.ErrorResponseWithErr(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusNoContent, nil)
}

// ModifyUser godoc
// @Summary      Modify user
// @Description  Updates user information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "Updated user data"
// @Success      200   {object}  map[string]models.User  "Updated user data"
// @Failure      400   {object}  utils.HTTPError        "Invalid request format"
// @Failure      500   {object}  utils.HTTPError        "Internal server error"
// @Router       /users/modify [post]
func (c UserController) ModifyUser(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := c.service.ModifyUser(context.Request.Context(), &user); err != nil {
		utils.ErrorResponseWithErr(context, http.StatusInternalServerError, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": user})
}

// ModifyUserLocation godoc
// @Summary      Modify user location
// @Description  Updates the location of a specific user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id        path      int         true  "User ID"
// @Param        location  body      models.LocationModifyRequest  true  "User with updated location"
// @Success      200       {object}  nil          "Location updated successfully"
// @Failure      400       {object}  utils.HTTPError  "Invalid user ID format or request"
// @Failure      500       {object}  utils.HTTPError  "Internal server error"
// @Router       /users/{id}/location [put]
func (c UserController) ModifyUserLocation(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	var user models.LocationModifyRequest
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "Invalid user ID format")
		return
	}
	if err := context.ShouldBindJSON(&user); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "Invalid request format")
		return
	}
	if err := c.service.ModifyLocation(context, id, user.Location); err != nil {
		utils.ErrorResponseWithErr(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, nil)
}

// BlockUserById godoc
// @Summary      Block user
// @Description  Blocks a user by ID
// @Tags         Users
// @Accept       json
// @Produce      plain
// @Param        id   path      int  true  "User ID"
// @Success      200  {string}  string  "User blocked successfully"
// @Failure      400  {object}  utils.HTTPError  "Invalid user ID format"
// @Failure      500  {object}  utils.HTTPError  "Internal server error"
// @Router       /users/block/{id} [put]
func (c UserController) BlockUserById(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// TODO: Add reason and blockerId
	if err := c.service.BlockUser(context.Request.Context(), id, "", nil, nil); err != nil {
		utils.ErrorResponseWithErr(context, http.StatusInternalServerError, err)
		return
	}

	context.String(http.StatusOK, "User blocked successfully")
}

// ModifyUserPasssword godoc
// @Summary      Modify user password
// @Description  Updates the password of a specific user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id        path      int         true  "User ID"
// @Param        password  body      models.PasswordModifyRequest  true  "User with updated password"
// @Success      200       {object}  nil          "Pasword updated successfully"
// @Failure      400       {object}  utils.HTTPError  "Invalid user ID format or request"
// @Failure      500       {object}  utils.HTTPError  "Internal server error"
// @Router       /users/{id}/password [put]
func (c UserController) ModifyUserPasssword(ctx *gin.Context) {
	var id, err = strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	var user models.PasswordModifyRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusBadRequest, err)
		return
	}
	if c.service.ModifyPassword(ctx.Request.Context(), id, user.Password) != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nil)

}

func (ct UserController) ValidateToken(c *gin.Context) {
	c.JSON(http.StatusOK, c.Request.Context().Value("jwt_info"))
}
