package server

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	chatroom string
	Egress   chan Event // helps avoid concurrent writes to the websocket connection
}

func (client *Client) readMessages() {
	defer func() {
		client.Hub.removeClient(client)
	}()
	// Setting pong-wait deadlines
	if err := client.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("pong timeout: %v", err)
		return
	}
	client.Conn.SetReadLimit(ClientReadLimit)
	client.Conn.SetPongHandler(client.PongHandler)

	for {
		_, payload, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("failed to unmarshal payload: %v", err)
			break
		}

		if err := client.Hub.checkEvent(request, client); err != nil {
			log.Printf("invalid event: %v", err)
		}
	}
}

func (client *Client) writeMessages() {
	defer func() {
		client.Hub.removeClient(client)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-client.Egress:
			if !ok {
				if err := client.Conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("failed to marshal message: %v", err)
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("failed to send message: %v", err)
			}
			log.Println("message sent")

		case <-ticker.C:
			log.Println("ping client")
			// Ping-ing the client
			if err := client.Conn.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Printf("writeMessage error: %v", err)
				return
			}
		}
	}
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{Conn: conn, Hub: hub, Egress: make(chan Event)}
}
