package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Load .env only in local development
	// if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
	// 	_ = godotenv.Load()
	// }
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("No .env file found")
	// 	log.Fatal(err)
	// }
	// services.InitJWTService()
	// client, err := db.InitDB()
	// if err != nil {
	// 	log.Fatal("Failed to initialize database:", err)
	// }
	// // initialize services with mongo client before starting the HTTP server
	// services.InitDB(client)

	// // pass mongo client to server and start it
	// newServer := server.NewWebServer(client)
	// if err := newServer.Connect(); err != nil {
	// 	log.Fatal(err)
	// }
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Railway OK"))
	})

	log.Println("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
