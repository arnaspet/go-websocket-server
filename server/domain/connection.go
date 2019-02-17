package domain

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ConnectionHolder interface {
	getConnection() *websocket.Conn
	receiveMessage(message []byte)
	getId() uint
}

type Connection struct {
	websocketConn *websocket.Conn
	logger        *logrus.Logger
	id            uint
	queue         chan []byte
}

func (s *Connection) getConnection() *websocket.Conn {
	return s.websocketConn
}

func (s *Connection) getId() uint {
	return s.id
}

func (s *Connection) receiveMessage(message []byte) {
	s.queue <- message
}

func (s *Connection) initMessageSender() {
	go func() {
		for {
			select {
			case msg := <-s.queue:
				s.logger.Debugf("Sending message to subscriber %v: %s", s.id, msg)
				err := s.websocketConn.WriteMessage(websocket.TextMessage, msg)

				if err != nil {
					s.logger.Error("write: ", err)
				}
			}
		}
	}()
}
