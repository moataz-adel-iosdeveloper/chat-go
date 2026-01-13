package ChatControllers

import (
	"chat-go/internal/models"
	conversationsRepository "chat-go/internal/repositories/conversations"
	"log"
	"sync"
)

type Room struct {
	clients map[*Client]bool
}

type ChatManager struct {
	rooms map[string]Room // Map of conversationID to room clients
	sync.RWMutex
}

func NewChatManager() *ChatManager {
	return &ChatManager{
		rooms: make(map[string]Room),
	}
}

func (m *ChatManager) RegisterClient(chatClient *Client, conversation *models.Conversation) {
	m.Lock()
	defer m.Unlock()
	conversationID := conversation.ID.Hex()

	// Create room if it doesn't exist
	if _, ok := m.rooms[conversationID]; !ok {
		m.rooms[conversationID] = Room{
			clients: make(map[*Client]bool),
		}
	}

	// Add client to the room
	room := m.rooms[conversationID]
	room.clients[chatClient] = true
	m.rooms[conversationID] = room
}

func (m *ChatManager) UnregisterClient(chatClient *Client, conversationID string) {
	m.Lock()
	defer m.Unlock()

	room, exists := m.rooms[conversationID]
	if !exists {
		return
	}
	// close connection to client
	chatClient.connection.Close()
	// delete client from map
	delete(room.clients, chatClient)

	// If room is empty, delete it
	if len(room.clients) == 0 {
		delete(m.rooms, conversationID)
	} else {
		m.rooms[conversationID] = room
	}
}

// return message to one client with message or error message
func (m *ChatManager) sentMessageToOneClient(to *Client, event SocketEvent) {
	m.RLock()
	defer m.RUnlock()
	to.writeChannel <- event
}

// BroadcastToRoom sends a message to all clients in a specific conversation/room
func (m *ChatManager) BroadcastToRoom(conversationID string, event SocketEvent, sender *Client) {
	m.RLock()
	room, exists := m.rooms[conversationID]
	m.RUnlock()

	if !exists {
		return
	}

	m.RLock()
	defer m.RUnlock()
	for client := range room.clients {
		// Optionally, skip sending the message back to the sender
		if client != sender {
			client.writeChannel <- event
		}
	}
}

func (m *ChatManager) saveMessageOnDatabase(conversationID string, message []byte, sender *Client) error {
	// get conversation by ID
	conversation, err := conversationsRepository.FindConversationByID(conversationID)
	if err != nil || conversation == nil {
		log.Println("Error finding conversation:", err)
		return err
	}
	// write message to messages collection database
	insertedMessage, err := sender.WriteMessageToDatabase(message, conversation)
	if err != nil {
		log.Println("Error writing message to database:", err)
		return err
	}
	// update conversation collection with last message
	err = sender.updateConnectionLastMessage(conversation, insertedMessage)
	if err != nil {
		log.Println("Error updating conversation last message:", err)
		return err
	}
	return nil
}

// // GetRoomClientCount returns the number of clients in a specific room
func (m *ChatManager) GetRoomClientCount(conversationID string) int {
	m.RLock()
	defer m.RUnlock()

	if room, exists := m.rooms[conversationID]; exists {
		return len(room.clients)
	}
	return 0
}
