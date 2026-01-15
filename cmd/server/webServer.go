package server

import (
	"chat-go/internal/middlewares"
	"chat-go/internal/routes"
	"os"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	Address string
	DB      *mongo.Client
}

func NewWebServer(client *mongo.Client) *server {
	appPort := os.Getenv("PORT")
	if appPort == "" {
		log.Println("No API_ADDRESS found in .env, using default :8080")
		appPort = "8080"
	}
	return &server{
		Address: ":" + appPort,
		DB:      client,
	}
}

// Start runs the HTTP server for the API.
func (webServer *server) Connect() error {
	router := mux.NewRouter()
	// global middlewares
	router.Use(middlewares.LoggerMiddleware)

	// ======================
	// Static files (CSS / JS / Images)
	// ======================
	router.PathPrefix("/assets/").
		Handler(http.StripPrefix(
			"/assets/",
			http.FileServer(http.Dir("./assets")),
		))

	// rotue to web api server
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	routes.NewRoutes().APIRoutes(subrouter)

	// Web Views Routes
	viewRouter := router.NewRoute().Subrouter()
	routes.NewRoutes().WebRoutes(viewRouter)

	// WebSocket Routes
	wsRouter := router.PathPrefix("/ws").Subrouter()
	routes.NewRoutes().SocketRoutes(wsRouter)

	log.Printf("starting server at %s", webServer.Address)
	return http.ListenAndServe(webServer.Address, middlewares.EnableCORS(router))
}
