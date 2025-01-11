package server

import "sync"

type ClientList map[*Client]bool
type Hub struct {
	clients ClientList
	sync.RWMutex
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
	return &Hub{clients: make(ClientList)}
}
