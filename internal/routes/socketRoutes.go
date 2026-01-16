package routes

import (
	ChatControllers "chat-go/internal/controllers/chatControllers"
	"chat-go/internal/middlewares"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (r *Route) SocketRoutes(router *mux.Router) {
	log.Println("Server started from Socket Routes")
	// Define WebSocket related routes here if needed
	router.Handle("/chat", middlewares.AuthMiddlewareSocket(http.HandlerFunc(ChatControllers.ChatHandler)))
}
