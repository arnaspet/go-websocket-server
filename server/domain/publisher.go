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

func(p *Publisher) broadcastMessage(message []byte) {
	messageToSend := ReplaceBytes(message)

	sendMessage := func(subscriber *Subscriber, message []byte) {
		p.logger.Debugf("Publishing message to subscriber %v", subscriber.id)
		err := subscriber.getConnection().WriteMessage(websocket.TextMessage, message)

		if err != nil {
			p.logger.Error("write: ", err)
		}
	}

	for subscriber := range p.pool.subscribers {
		go sendMessage(subscriber, messageToSend)
	}
}

func (p *Publisher) initMessageHandler() {
	for {
		mt, message, err := p.conn.ReadMessage()

		if err != nil {
			p.logger.Error("read: ", err)
			break
		}

		if mt != websocket.TextMessage {
			p.pool.closeWebsocketConnection(p)
			break
		}

		p.logger.Debugf("Message from publisher #%d received: %s", p.id, message)
		p.broadcastMessage(message)
	}
}
