package server

import (
	"encoding/json"
	"time"
)

const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
)

type EventHandler func(event Event, c *Client) error

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"` // Not constraining the payload by assigning a type
}

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}
