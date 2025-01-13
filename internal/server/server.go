package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	AllowedOrigin   = "http://localhost:3001"
	ClientReadLimit = 512 // Precaution for Jumbo-Frames
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     checkOrigin,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(hub *Hub, res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println("failed to upgrade connection: ", err)
		return
	}
	client := NewClient(conn, hub)
	hub.addClient(client)

	// Client Processes
	go client.readMessages()
	go client.writeMessages()

	conn.Close()
}

func checkOrigin(req *http.Request) bool {
	origin := req.Header.Get("Origin")
	switch origin {
	case AllowedOrigin:
		return true
	default:
		return false
	}
}
