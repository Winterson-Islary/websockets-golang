package server

import (
	"fmt"
	"log"
	"time"
)

const (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

func (client *Client) PongHandler(pongMsg string) error {
	log.Println("pong from client")
	if err := client.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return fmt.Errorf("could not update read deadline: %w", err)
	}
	return nil
}
