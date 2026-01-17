package controllers

import (
	"chat-go/internal/models"
	conversationsRepository "chat-go/internal/repositories/conversations"
	MessagesRepository "chat-go/internal/repositories/messages"
	userRepository "chat-go/internal/repositories/user"
	"errors"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if conversationResponse == nil {
		conversationResponse = []models.ConversationResponse{}
	}
	return models.SuccessResponse("users retrieved successfully", conversationResponse), http.StatusOK
}

func GetAllMessages(req *http.Request) (response models.APIResponse, statusCode int) {

	// Get the authenticated user ID from context
	authUserID := req.Context().Value("userID")
	if authUserID == nil {
		log.Println("Error parsing auth ID")
		return models.ErrorResponse("Error parsing auth ID", nil), http.StatusBadRequest
	}

	// get user from auth user id
	authUser, err := userRepository.FindUserByID(authUserID.(string))
	if err != nil {
		log.Println("Error get auth user")
		return models.ErrorResponse("Error get auth user", nil), http.StatusBadRequest
	}

	var conversation *models.Conversation

	// get another user id from query params
	conversationID := req.URL.Query().Get("conversation_id")
	userID := req.URL.Query().Get("user_id")

	if conversationID == "" {
		// conversation ID is empaty it mean create new Conversations or search for Conversations it by two users
		if userID == "" {
			log.Println("Error parsing conversation ID")
			return models.ErrorResponse("Error parsing conversation ID", nil), http.StatusBadRequest
		}
		// get or create conversation
		conversation, err = conversationsRepository.FindOrCreateConversationsByTwoUsersID(userID, authUserID.(string))
		if err != nil {
			log.Println("Error on create or find conversation")
			return models.ErrorResponse("Error create or find conversation", nil), http.StatusBadRequest
		}
		userID, err = GetOtherUserID(conversation.ParticipantIDs, authUserID.(string))
		if err != nil {
			log.Println("Error on get other conversation user")
			return models.ErrorResponse("Error create or find conversation", nil), http.StatusBadRequest
		}
		// finalConversation = conversation
		conversationID = conversation.ID.Hex()
	} else {
		conversation, err = conversationsRepository.FindConversationByID(conversationID)
		if err != nil {
			log.Println("Error find conversation", err)
			return models.ErrorResponse("Error find conversation", nil), http.StatusBadRequest
		}
		userID, err = GetOtherUserID(conversation.ParticipantIDs, authUserID.(string))
		if err != nil {
			log.Println("Error on get other conversation user")
			return models.ErrorResponse("Error create or find conversation", nil), http.StatusBadRequest
		}
	}

	log.Println("other user id raw:", userID)
	// get other user from user id
	otherUser, err := userRepository.FindUserByID(userID)
	if err != nil {
		log.Println("Error get other user", err)
		return models.ErrorResponse("Error get other user", nil), http.StatusBadRequest
	}

	// Get all messages by ConversationID
	messages, err := MessagesRepository.FindAllMessagesByConversationID(conversationID)
	if err != nil {
		log.Println("Error by get messages from data", err)
		return models.ErrorResponse("Error on get messages", nil), http.StatusBadRequest
	}

	// create responce
	responseModel := models.AllMessagesResponce{
		Messages:     messages,
		User:         authUser,
		OtherUser:    otherUser,
		Conversation: conversation,
	}

	return models.SuccessResponse("messages retrieved successfully", responseModel), http.StatusOK
}

func GetOtherUserID(participants []primitive.ObjectID, authID string) (string, error) {
	for _, id := range participants {
		if id.Hex() != authID {
			return id.Hex(), nil
		}
	}
	return "", errors.New("other user not found")
}
