package domain

import (
	"github.com/gorilla/websocket"
)

type Publisher struct {
	conn *websocket.Conn
	id uint
}
