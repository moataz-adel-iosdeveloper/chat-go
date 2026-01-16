package controllers

import (
	"chat-go/internal/models"
	conversationsRepository "chat-go/internal/repositories/conversations"
	MessagesRepository "chat-go/internal/repositories/messages"
	"log"
	"net/http"
)

func GetAllConversations(req *http.Request) (response models.APIResponse, statusCode int) {

	// Get the authenticated user ID from context
	authUserID := req.Context().Value("userID")
	if authUserID == nil {
		log.Println("Error parsing auth ID")
		return models.ErrorResponse("Error parsing auth ID", nil), http.StatusBadRequest
	}

	// Get all Conversations by userID
	var conversations []*models.Conversation
	conversations, err := conversationsRepository.FindConversationsByUserID(authUserID.(string))
	if err != nil {
		log.Println("Error by get conversations from data", err)
		return models.ErrorResponse("Error on get conversations", nil), http.StatusBadRequest
	}

	var conversationResponse []models.ConversationResponse
	for _, con := range conversations {
		participantIDs := make([]string, 0, len(con.ParticipantIDs))
		for _, id := range con.ParticipantIDs {
			participantIDs = append(participantIDs, id.Hex())
		}
		conversationResponse = append(conversationResponse, models.ConversationResponse{
			ID:             con.ID.Hex(),
			ParticipantIDs: participantIDs,
			LastMessage:    con.LastMessage,
			CreatedAt:      con.CreatedAt.Time(),
			UpdatedAt:      con.UpdatedAt.Time(),
		})
	}
	return models.SuccessResponse("users retrieved successfully", conversationResponse), http.StatusOK
}

func GetAllMessages(req *http.Request) (response models.APIResponse, statusCode int) {
	// get another user id from query params
	conversationID := req.URL.Query().Get("conversation_id")
	if conversationID == "" {
		log.Println("Error parsing conversation ID")
		return models.ErrorResponse("Error parsing conversation ID", nil), http.StatusBadRequest
	}

	// Get all messages by ConversationID
	messages, err := MessagesRepository.FindAllMessagesByConversationID(conversationID)
	if err != nil {
		log.Println("Error by get messages from data", err)
		return models.ErrorResponse("Error on get messages", nil), http.StatusBadRequest
	}

	return models.SuccessResponse("messages retrieved successfully", messages), http.StatusOK
}
