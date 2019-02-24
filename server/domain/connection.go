package domain

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ConnectionHolder interface {
	getConnection() *websocket.Conn
	receiveMessage(message *Message)
	getId() uint64
}

type Message struct {
	Content []byte
	MsgType int
}

type Connection struct {
	websocketConn *websocket.Conn
	logger        *logrus.Logger
	id            uint64
	queue         chan *Message
}

func NewConnection(websocketConn *websocket.Conn, logger *logrus.Logger, id uint64, queue chan *Message) *Connection {
	connection := &Connection{websocketConn: websocketConn, logger: logger, id: id, queue: queue}
	connection.initMessageSender()

	return connection
}

func (s *Connection) getConnection() *websocket.Conn {
	return s.websocketConn
}

func (s *Connection) getId() uint64 {
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
				s.logger.Debugf("Sending message to Connection %v: %s", s.id, msg)
				err := s.websocketConn.WriteMessage(msg.MsgType, msg.Content)

				if err != nil {
					s.logger.Error("write: ", err)
				}
			}
		}
	}()
}
