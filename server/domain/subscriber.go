package domain

import (
	"github.com/gorilla/websocket"
)

type Subscriber struct {
	pool *ConnectionPool
	conn *Connection
}

func (s *Subscriber) getConnection() *websocket.Conn {
	return s.conn.getConnection()
}

func (s *Subscriber) receiveMessage(message *Message) {
	s.conn.receiveMessage(message)
}

func (s *Subscriber) getId() uint {
	return s.conn.getId()
}

func (s *Subscriber) initMessageHandler() {
	for {
		mt, message, err := s.getConnection().ReadMessage()

		if mt != websocket.TextMessage {
			s.pool.closeWebsocketConnection(s)
			break
		}

		if err != nil {
			s.conn.logger.Error("read: ", err)
			break
		}
		s.conn.logger.Debugf("recv: %s", message)

		err = s.getConnection().WriteMessage(mt, ReplaceBytes(message))

		if err != nil {
			s.conn.logger.Error("write: ", err)
			break
		}
	}
}
