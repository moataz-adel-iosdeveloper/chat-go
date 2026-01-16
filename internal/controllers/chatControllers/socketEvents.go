package ChatControllers

import (
	"encoding/json"
	"log"
)

type SocketEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type TextMessagePayload struct {
	Message string `json:"message"`
}

type ErrorPayload struct {
	Code    string `json:"error"`
	Message string `json:"message"`
}

type EventHandler func(event SocketEvent, client *Client) error

const (
	TextMessageEvent    = "text_message"
	ErrorEvent          = "error"
	InvaledPayloadError = "INVALID_PAYLOAD"
	ErrorParsingPayload = "ERROR_PARSING"
)

var eventHandler = map[string]EventHandler{
	TextMessageEvent: handleTextMessage,
}

func handleTextMessage(event SocketEvent, client *Client) error {
	var message TextMessagePayload
	if err := json.Unmarshal(event.Payload, &message); err != nil {
		event := NewErrorEvent(ErrorParsingPayload, "error on parsing payload json")
		client.manager.sentMessageToOneClient(client, event)
		return err
	}
	log.Printf("Client %s sent message: %s", client.userID, message.Message)
	// Broadcast message to all clients in the same conversation/room
	client.manager.BroadcastToRoom(client.conversationID, event, client)
	// save message to database
	client.manager.saveMessageOnDatabase(client.conversationID, []byte(message.Message), client)
	return nil
}

func NewErrorEvent(code, message string) SocketEvent {
	payload, _ := json.Marshal(ErrorPayload{
		Code:    code,
		Message: message,
	})

	return SocketEvent{
		Type:    ErrorEvent,
		Payload: payload,
	}
}
