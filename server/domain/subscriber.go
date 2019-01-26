package domain

import "github.com/gorilla/websocket"

type Subscriber struct {
	conn *websocket.Conn
	id uint
}
