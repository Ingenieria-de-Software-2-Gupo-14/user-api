package models

type ChatMessageRequest struct {
	Message string `json:"message" binding:"required"`
}

type ChatRatingRequest struct {
	Rating int `json:"rating" binding:"required,gte=1,lte=5"`
}

type ChatFeedbackRequest struct {
	Feedback string `json:"feedback" binding:"required"`
}

type ChatMessage struct {
	MessageId int    `json:"message_id"`
	UserId    int    `json:"user_id"`
	Message   string `json:"message"`
	Sender    string `json:"sender"`
	TimeSent  string `json:"time_sent"`
	Rating    int    `json:"rating"`
	Feedback  string `json:"feedback"`
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}
