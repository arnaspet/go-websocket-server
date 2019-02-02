package domain

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Subscriber struct {
	pool *ConnectionPool
	conn *websocket.Conn
	logger *logrus.Logger
	id uint
}


func (s *Subscriber) getConnection() *websocket.Conn {
	return s.conn
}

func (s *Subscriber) getId() uint {
	return s.id
}

func (s *Subscriber) initMessageHandler() {
	for {
		mt, message, err := s.conn.ReadMessage()

		if mt != websocket.TextMessage {
			s.pool.closeWebsocketConnection(s)
			break
		}

		if err != nil {
			s.logger.Error("read: ", err)
			break
		}
		s.logger.Debugf("recv: %s", message)

		err = s.conn.WriteMessage(mt, ReplaceBytes(message))

		if err != nil {
			s.logger.Error("write: ", err)
			break
		}
	}
}
