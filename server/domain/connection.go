package domain

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ConnectionHolder interface {
	getConnection() *websocket.Conn
	receiveMessage(message *Message)
	getId() uint
}

type Message struct {
	content []byte
	msgType int
}

type Connection struct {
	websocketConn *websocket.Conn
	logger        *logrus.Logger
	id            uint
	queue         chan *Message
}

func (s *Connection) getConnection() *websocket.Conn {
	return s.websocketConn
}

func (s *Connection) getId() uint {
	return s.id
}

func (s *Connection) receiveMessage(message *Message) {
	s.queue <- message
}

func (s *Connection) initMessageSender() {
	go func() {
		for {
			select {
			case msg := <-s.queue:
				s.logger.Debugf("Sending message to subscriber %v: %s", s.id, msg)
				err := s.websocketConn.WriteMessage(msg.msgType, msg.content)

				if err != nil {
					s.logger.Error("write: ", err)
				}
			}
		}
	}()
}
