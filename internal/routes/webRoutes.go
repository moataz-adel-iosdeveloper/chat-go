package routes

import (
	"chat-go/internal/middlewares"
	"chat-go/internal/templates/front"
	"net/http"

	"github.com/gorilla/mux"
)

func (r *Route) WebRoutes(router *mux.Router) {
	router.HandleFunc("/login", r.loginPage).Methods("GET")
	router.HandleFunc("/register", r.registerPage).Methods("GET")

	// Protected routes
	protected := router.NewRoute().Subrouter()
	protected.Use(middlewares.WebAuthMiddleware(r.Session))

	protected.HandleFunc("/", r.HomePage).Methods("GET")
	protected.HandleFunc("/home", r.HomePage).Methods("GET")
	protected.HandleFunc("/chat", r.ChatPage).Methods("GET")
}

func (r *Route) HomePage(w http.ResponseWriter, req *http.Request) {
	// Get current user ID from context/session
	userID := getUserIDFromContext(req.Context())

	component := front.Home(userID)
	component.Render(req.Context(), w)
}

func (r *Route) loginPage(w http.ResponseWriter, req *http.Request) {
	component := front.Login(false, "")
	component.Render(req.Context(), w)
}

func (r *Route) registerPage(w http.ResponseWriter, req *http.Request) {
	component := front.Register(false, "")
	component.Render(req.Context(), w)
}

func (r *Route) ChatPage(w http.ResponseWriter, req *http.Request) {
	// Get conversation ID from query params
	conversationID := req.URL.Query().Get("conversation_id")
	if conversationID == "" {
		conversationID = req.URL.Query().Get("user_id")
	}

	// Get other user name (simplified - in real app, fetch from DB)
	otherUserName := "Chat"

	data := front.ChatPageData{
		ConversationID: conversationID,
		Messages:       []front.ChatMessage{},
		OtherUserName:  otherUserName,
	}

	component := front.ChatRoom(data)
	component.Render(req.Context(), w)
}

func getUserIDFromContext(ctx interface{}) string {
	if ctx == nil {
		return ""
	}
	if userID, ok := ctx.(string); ok {
		return userID
	}
	return ""
}
