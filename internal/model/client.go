package model

import "github.com/gorilla/websocket"

type Client struct {
	Conn     *websocket.Conn
	Send     chan []byte
	ClientId string
}

type Message struct {
	ClientID string `json:"clientId"`
	Content  string `json:"content"`
}
