package routes

import (
	"chat-go/internal/config"
	"chat-go/internal/middlewares"
	"chat-go/internal/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Session config.Store
}

func NewRoutes() *Route {
	var store = config.NewCookieStore([]byte("super-secret-key"))
	return &Route{
		Session: store,
	}
}

func (r *Route) APIRoutes(router *mux.Router) {
	router.HandleFunc("/login", r.loginHandler).Methods("POST")
	router.HandleFunc("/register", r.registerHandler).Methods("POST")
	router.Handle("/allUsers", middlewares.
		AuthMiddleware(http.HandlerFunc(r.allUserHandler))).Methods("GET")
}

func (r *Route) WriterResponse(w http.ResponseWriter, req *http.Request, response models.APIResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
