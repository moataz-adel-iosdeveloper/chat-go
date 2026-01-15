package config

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

func SocketServe(w http.ResponseWriter, req *http.Request) (*websocket.Conn, error) {
	log.Println("new connection on WebSocket")
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return nil, err
	}
	return conn, nil
}
