package routes

import (
	ChatControllers "chat-go/internal/controllers/chatController"
	"chat-go/internal/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func (r *Route) SocketRoutes(router *mux.Router) {
	// Define WebSocket related routes here if needed
	router.Handle("/chat", middlewares.AuthMiddlewareSocket(http.HandlerFunc(ChatControllers.ChatHandler)))
}
