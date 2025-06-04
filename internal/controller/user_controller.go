package controller

import (
	"net/http"
	"os"
	"strconv"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	services "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/utils"

	"github.com/gin-gonic/gin"
)

// UserController struct that contains a database with users
type UserController struct {
	service     services.UserService
	ruleService services.RulesService
}

// CreateController creates a controller
func CreateController(service services.UserService, ruleService services.RulesService) *UserController {
	return &UserController{service: service, ruleService: ruleService}
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
// @Param        id    path      int         true  "User ID"
// @Param        user  body      models.UserUpdateDto  true  "Updated user data"
// @Success      200   {object}  map[string]models.UserUpdateDto  "Updated user data"
// @Failure      400   {object}  utils.HTTPError        "Invalid user ID format or request format"
// @Failure      500   {object}  utils.HTTPError        "Internal server error"
// @Router       /users/{id} [put]
func (c UserController) ModifyUser(context *gin.Context) {
	var id, err = strconv.Atoi(context.Param("id"))
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var user models.UserUpdateDto
	if err := context.ShouldBindJSON(&user); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := c.service.ModifyUser(context.Request.Context(), id, user); err != nil {
		utils.ErrorResponseWithErr(context, http.StatusInternalServerError, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": user})
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

// NotifyUsers godoc
// @Summary      Send a notification to users
// @Description  Send a notification to users sent in body
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Notification  body  models.NotifyRequest true  "Notification payload"
// @Success      200       {object}  nil          "Users notified successfully"
// @Failure      400       {object}  utils.HTTPError  "Invalid request"
// @Failure      500       {object}  utils.HTTPError  "Internal server error"
// @Router       /users/notify [post]
func (c UserController) NotifyUsers(ctx *gin.Context) {
	var notifyRequest models.NotifyRequest
	if err := ctx.ShouldBindJSON(&notifyRequest); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request format")
		return
	}
	cont := ctx.Request.Context()
	for _, userID := range notifyRequest.Users {
		user, err := c.service.GetUserById(cont, userID)
		if err != nil {
			if err == repositories.ErrNotFound {
				continue
			}
			utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
			continue
		}
		err = c.service.AddNotification(cont, userID, notifyRequest.NotificationText)
		if err != nil {
			utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
			continue
		}
		errMail := sendNotifByEmail(user.Email, notifyRequest.NotificationText)
		if errMail != nil {
			utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, errMail)
			continue
		}
	}
	ctx.JSON(http.StatusOK, nil)
}

func sendNotifByEmail(email string, text string) error {
	from := mail.NewEmail("ClassConnect service", "bmorseletto@fi.uba.ar")
	subject := "Verification Code"
	to := mail.NewEmail("User", email)
	content := mail.NewContent("text/plain", text)
	message := mail.NewV3MailInit(from, subject, to, content)

	client := sendgrid.NewSendClient(os.Getenv("EMAIL_API_KEY"))
	_, err := client.Send(message)
	return err
}

// GetUserNotifications godoc
// @Summary      Send a notification to users
// @Description  Send a notification to users sent in body
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id        path      int         true  "User ID"
// @Success      200       {object}  models.Notifications          "Users notified successfully"
// @Failure      500       {object}  utils.HTTPError  "Internal server error"
// @Router       /users/{id}/notifications [get]
func (c UserController) GetUserNotifications(ctx *gin.Context) {
	var id, err = strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	var notifs models.Notifications
	notifs, err = c.service.GetUserNotifications(ctx.Request.Context(), id)
	if err != nil {
		if err == repositories.ErrNotFound {
			utils.ErrorResponse(ctx, http.StatusNotFound, "User not found")
		} else {
			utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, notifs)
}

// AddRule godoc
// @Summary      Create a new rule
// @Description  Create a new rule
// @Tags         Rules
// @Accept       json
// @Produce      plain
// @Param        request body models.Rule true "Rule creation Details"
// @Success      201       {object}  nil          "Rule created correctly"
// @Failure      400  {object}  utils.HTTPError "Invalid request format"
// @Failure      500  {object}  utils.HTTPError "Internal server error"
// @Router       /rules [post]
func (c UserController) AddRule(ctx *gin.Context) {
	auth, _ := ctx.Cookie("Authorization")
	claims, err := models.ParseToken(auth)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	var rule models.Rule
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request format")
		return
	}
	err = c.ruleService.CreateRule(ctx, rule, userId)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, nil)
}

// DeleteRule godoc
// @Summary      Delete user by ID
// @Description  Removes a user from the database
// @Tags         Rules
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Rule ID"
// @Success      204  {object}  nil  "Rule successfully deleted"
// @Failure      400  {object}  utils.HTTPError  "Invalid user ID format"
// @Failure      500  {object}  utils.HTTPError  "Internal server error"
// @Router       /rules/{id} [delete]
func (c UserController) DeleteRule(ctx *gin.Context) {
	var id, err = strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	auth, _ := ctx.Cookie("Authorization")
	claims, err := models.ParseToken(auth)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	err = c.ruleService.DeleteRule(ctx, id, userId)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

// GetRules godoc
// @Summary      Get all rules
// @Description  Returns a list of all rules in the system
// @Tags         Rules
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string][]models.Rule  "List of rules"
// @Failure      500  {object}  utils.HTTPError          "Internal server error"
// @Router       /rules [get]
func (c UserController) GetRules(ctx *gin.Context) {
	rules, err := c.ruleService.GetRules(ctx)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": rules})
}

// GetAudits godoc
// @Summary      Get all audits
// @Description  Returns a list of all audits in the system
// @Tags         Rules
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string][]models.Audit  "List of audits"
// @Failure      500  {object}  utils.HTTPError          "Internal server error"
// @Router       /rules [get]
func (c UserController) GetAudits(ctx *gin.Context) {
	audits, err := c.ruleService.GetAudits(ctx)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": audits})
}

// ModifyRule godoc
// @Summary      Modify rule password
// @Description  Updates the contents of a rule
// @Tags         Rules
// @Accept       json
// @Param        id        path      int         true  "Rule ID"
// @Param        modifications  body      models.RuleModify  true  "Elements to modify"
// @Success      200       {object}  nil          "rule updated successfully"
// @Failure      400       {object}  utils.HTTPError  "Invalid user ID format or request"
// @Failure      500       {object}  utils.HTTPError  "Internal server error"
// @Router       /rules/{id} [put]
func (c UserController) ModifyRule(ctx *gin.Context) {
	var id, err = strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	auth, _ := ctx.Cookie("Authorization")
	claims, err := models.ParseToken(auth)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	var rule models.RuleModify
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request format")
		return
	}
	err = c.ruleService.ModifyRule(ctx.Request.Context(), id, rule, userId)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
