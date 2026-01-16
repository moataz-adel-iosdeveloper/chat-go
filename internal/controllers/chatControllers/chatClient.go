package ChatControllers

import (
	"chat-go/internal/models"
	conversationsRepository "chat-go/internal/repositories/conversations"
	messagesRepository "chat-go/internal/repositories/messages"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	connection     *websocket.Conn
	manager        *ChatManager
	writeChannel   chan SocketEvent // used to avoid concurince write on web socket connection
	conversationID string           // ID of the conversation/room this client belongs to
	userID         string           // ID of the user for this client
}

func NewClient(conn *websocket.Conn, manager *ChatManager, conversationID string, userID string) *Client {
	return &Client{
		connection:     conn,
		manager:        manager,
		writeChannel:   make(chan SocketEvent),
		conversationID: conversationID,
		userID:         userID,
	}
}

func (c *Client) ReadMessages() {
	defer c.manager.UnregisterClient(c, c.conversationID)
	for {
		messageType, message, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error read message: %v", err)
			}
			break
		}
		log.Printf("reseved new message with type:%v", messageType)
		var event SocketEvent
		if err := json.Unmarshal(message, &event); err != nil {
			event := NewErrorEvent(InvaledPayloadError, "invalid event:")
			c.manager.sentMessageToOneClient(c, event)
			log.Println("invalid event:", err)
			continue
		}
		if err := c.routeEvent(event, c); err != nil {
			event := NewErrorEvent(InvaledPayloadError, "invalid event:")
			c.manager.sentMessageToOneClient(c, event)
			log.Println("event handling error:", err)
			continue
		}

	}
}

func (c *Client) reuturnErrorEventToClient(event SocketEvent) {
	c.manager.sentMessageToOneClient(c, event)
}

func (c *Client) routeEvent(event SocketEvent, client *Client) error {
	handler, exists := eventHandler[event.Type]
	if !exists {
		return fmt.Errorf("unsupported event type: %s", event.Type)
	}

	return handler(event, client)
}

func (c *Client) WriteMessage() {
	defer c.manager.UnregisterClient(c, c.conversationID)

	for message := range c.writeChannel {
		data, err := json.Marshal(message)
		if err != nil {
			log.Println("failed to parse message", err)
			continue
		}

		if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("failed to send message to client", err)
			return
		}
		log.Println("message sent success")
	}

	// channel closed
	if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
		log.Println("lose connection with client", err)
	}
}

func (c *Client) WriteMessageToDatabase(message []byte, conversation *models.Conversation) (*models.Message, error) {
	// convert userID to ObjectID
	objID, err := primitive.ObjectIDFromHex(c.userID)
	if err != nil {
		return nil, err
	}
	// create message model
	newMessage := &models.Message{
		ConversationID: conversation.ID,
		SenderID:       objID,
		Content:        string(message),
		ContentType:    "text",
		IsRead:         false,
		CreatedAt:      primitive.NewDateTimeFromTime(time.Now()),
	}
	// insert message into database
	insertedMessage, err := messagesRepository.CreateMessage(newMessage)
	if err != nil {
		return nil, err
	}

	return insertedMessage, nil
}

func (c *Client) updateConnectionLastMessage(conversation *models.Conversation, message *models.Message) error {
	log.Printf("Writing message to client %s: %s", c.userID, message.Content)
	conversation.LastMessage = message
	conversation.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err := conversationsRepository.UpdateConversation(conversation)
	return err
}
