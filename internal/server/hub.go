package server

import (
	"errors"
	"fmt"
	"sync"
)

type ClientList map[*Client]bool

type Hub struct {
	clients ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

func (hub *Hub) addClient(client *Client) {
	hub.Lock()
	defer hub.Unlock()
	hub.clients[client] = true
}
func (hub *Hub) removeClient(client *Client) {
	hub.Lock()
	defer hub.Unlock()
	if _, ok := hub.clients[client]; ok {
		client.Conn.Close()
		delete(hub.clients, client)
	}
}

func NewHub() *Hub {
	hub := &Hub{clients: make(ClientList), handlers: make(map[string]EventHandler)}
	hub.setupEventHandlers()
	return hub
}

func (hub *Hub) setupEventHandlers() {
	hub.handlers[EventSendMessage] = SendMessage
}

func (hub *Hub) checkEvent(event Event, client *Client) error {
	if handler, ok := hub.handlers[event.Type]; ok {
		if err := handler(event, client); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("received invalid event type")
	}
}

// EVENT FUNCTIONS
func SendMessage(event Event, client *Client) error {
	fmt.Println(event)
	return nil
}
