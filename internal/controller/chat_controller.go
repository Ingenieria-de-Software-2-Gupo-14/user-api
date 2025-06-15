package controller

import (
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ChatController struct {
	chatService services.ChatService
}

func NewChatsController(chatService services.ChatService) *ChatController {
	return &ChatController{
		chatService,
	}
}

// SendMessage
// @Summary      Send a new message to the assitant
// @Description  Send a new message to the assitant
// @Tags         Chat
// @Accept       json
// @Produce      plain
// @Param        request body models.ChatMessageRequest true "User Registration Details"
// @Success      200  {object}  map[models.ChatMessage]int  "Ai assitant response"
// @Failure      400  {object}  utils.HTTPError "Invalid request format"
// @Failure      500  {object}  utils.HTTPError "Internal server error"
// @Router       /chat [post]
func (c ChatController) SendMessage(ctx *gin.Context) {
	tokenStr := GetAuthToken(ctx)
	claims, err := models.ParseToken(tokenStr)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	var message models.ChatMessageRequest
	if err := ctx.ShouldBindJSON(&message); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request format")
		return
	}
	err = c.chatService.NewUserMessage(ctx.Request.Context(), userId, message)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
	}
	aiMessage, err := c.chatService.SendToAi(ctx.Request.Context(), userId, message)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, gin.H{"data": aiMessage})
}

// GetMessages
// @Summary      Get recent messages
// @Description  Gets The messages of the last 2 days
// @Tags         Chat
// @Accept       json
// @Produce      plain
// @Success      200  {object}  []models.ChatMessage  "messages"
// @Failure      500  {object}  utils.HTTPError "Internal server error"
// @Router       /chat [get]
func (c ChatController) GetMessages(ctx *gin.Context) {
	tokenStr := GetAuthToken(ctx)
	claims, err := models.ParseToken(tokenStr)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	messages, err := c.chatService.GetMessages(ctx.Request.Context(), userId)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, messages)
}

// RateMessage
// @Summary      Add rating to a message
// @Description  Add rating to a message beteween 1 and 5
// @Tags         Chat
// @Accept       json
// @Produce      plain
// @Param        request body models.ChatRatingRequest true "Rating"
// @Param        message_id   path      int  true  "Message Id"
// @Success      200  {object}  nil
// @Failure      500  {object}  utils.HTTPError "Internal server error"
// @Router       /chat/{message_id}/rate [put]
func (c ChatController) RateMessage(ctx *gin.Context) {
	tokenStr := GetAuthToken(ctx)
	claims, err := models.ParseToken(tokenStr)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	messageId, err := strconv.Atoi(ctx.Param("message_id"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid message ID format")
		return
	}
	var rating models.ChatRatingRequest
	if err := ctx.ShouldBindJSON(&rating); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request format")
		return
	}
	err = c.chatService.UpdateMessageRating(ctx.Request.Context(), userId, messageId, rating.Rating)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, nil)
}

// FeedbackMessage
// @Summary      Add feedback to a message
// @Description  Add feedback to a message so the ai can better generate a response
// @Tags         Chat
// @Accept       json
// @Produce      plain
// @Param        request body models.ChatFeedbackRequest true "Rating"
// @Param        message_id   path      int  true  "Message Id"
// @Success      200  {object}  nil
// @Failure      500  {object}  utils.HTTPError "Internal server error"
// @Router       /chat/{message_id}/feedback [put]
func (c ChatController) FeedbackMessage(ctx *gin.Context) {
	tokenStr := GetAuthToken(ctx)
	claims, err := models.ParseToken(tokenStr)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
		return
	}
	messageId, err := strconv.Atoi(ctx.Param("message_id"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid message ID format")
		return
	}
	var feedback models.ChatFeedbackRequest
	if err := ctx.ShouldBindJSON(&feedback); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request format")
		return
	}
	err = c.chatService.UpdateMessageFeedback(ctx.Request.Context(), userId, messageId, feedback.Feedback)
	if err != nil {
		utils.ErrorResponseWithErr(ctx, http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, nil)
}
