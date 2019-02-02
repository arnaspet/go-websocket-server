package domain

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Publisher struct {
	pool *ConnectionPool
	conn *websocket.Conn
	logger *logrus.Logger
	id uint
}

func (p *Publisher) getConnection() *websocket.Conn {
	return p.conn
}

func (p *Publisher) getId() uint {
	return p.id
}

func(p *Publisher) BroadcastMessage(to []ConnectionHolder, message []byte) {
	for _, receiver := range to {
		receiver.getConnection().WriteMessage(websocket.TextMessage, message)
	}
}

func (p *Publisher) initMessageHandler() {
	for {
		mt, message, err := p.conn.ReadMessage()

		if mt != websocket.TextMessage {
			p.pool.closeWebsocketConnection(p)
			break
		}

		if err != nil {
			p.logger.Error("read: ", err)
			break
		}
		p.logger.Debugf("recv: %s", message)

		for _, subscriber := range p.pool.subscribers {
			p.logger.Printf("Publishing message to subscriber %v", subscriber.id)
			err = subscriber.getConnection().WriteMessage(mt, ReplaceBytes(message))

			if err != nil {
				p.logger.Error("write: ", err)
				break
			}
		}
	}
}
