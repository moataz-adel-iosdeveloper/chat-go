package ChatControllers

import (
	"chat-go/internal/config"
	"chat-go/internal/models"
	conversationsRepository "chat-go/internal/repositories/conversations"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var chatManager *ChatManager

func init() {
	// Initialize any required resources or configurations here
	chatManager = NewChatManager()
}

func ChatHandler(w http.ResponseWriter, req *http.Request) {

	// get another user id from query params
	userID := req.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	// Get the authenticated user ID from context
	authUserID := req.Context().Value("userID")
	if authUserID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if authUserID.(string) == userID {
		http.Error(w, "Cannot chat with yourself", http.StatusBadRequest)
		return
	}

	log.Printf("Authenticated user ID: %v, Chatting with user ID: %s", authUserID, userID)

	// Find a conversation between the two users
	conversation, err := conversationsRepository.FindConversationByTwoUserID(authUserID.(string), userID)
	if err != nil {
		log.Println("Error finding conversation:", err)
		http.Error(w, "Error finding conversation", http.StatusInternalServerError)
		return
	}

	// create a conversation between the two users
	if conversation == nil {
		senderObjID, err := primitive.ObjectIDFromHex(authUserID.(string))
		if err != nil {
			http.Error(w, "Error finding or Sender Id", http.StatusInternalServerError)
			return
		}
		reseverObjID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			http.Error(w, "Error finding or Resever Id", http.StatusInternalServerError)
			return
		}
		newConversation := &models.Conversation{
			ParticipantIDs: []primitive.ObjectID{senderObjID, reseverObjID},
			CreatedAt:      primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt:      primitive.NewDateTimeFromTime(time.Now()),
		}
		conversation, err = conversationsRepository.CreateConversation(newConversation)
		if err != nil {
			log.Println("Error creating conversation:", err)
			http.Error(w, "Error creating conversation", http.StatusInternalServerError)
			return
		}
	}

	connectionHandler(conversation, authUserID.(string), w, req)
}

func connectionHandler(conversation *models.Conversation, userID string, w http.ResponseWriter, req *http.Request) {
	connection, err := config.SocketServe(w, req)
	if err != nil {
		log.Println("WebSocket connection error:", err)
		return
	}
	client := NewClient(connection, chatManager, conversation.ID.Hex(), userID)
	chatManager.RegisterClient(client, conversation)
	go client.ReadMessages()
	go client.WriteMessage()
}
