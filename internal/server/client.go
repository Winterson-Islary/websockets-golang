package server

import "github.com/gorilla/websocket"

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Msg  chan []byte
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{Conn: conn, Hub: hub}
}
