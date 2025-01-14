package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
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
	hub.handlers[EventChangeRoom] = ChatRoomHandler
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
	var chatEvent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		return fmt.Errorf("invalid payload in request: %y", err)
	}
	var broadMessage NewMessageEvent
	broadMessage.Sent = time.Now()
	broadMessage.Message = chatEvent.Message
	broadMessage.From = chatEvent.From

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %y", err)
	}
	outgoingEvent := Event{
		Type:    EventNewMessage,
		Payload: data,
	}
	for _client := range client.Hub.clients {
		if _client.chatroom == client.chatroom {
			_client.Egress <- outgoingEvent
		}
	}
	return nil
}

func ChatRoomHandler(event Event, client *Client) error {
	var changeRoomEvent ChangeRoomEvent
	if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
		return fmt.Errorf("bad payload request: %v", err)
	}
	client.chatroom = changeRoomEvent.Name
	return nil
}
