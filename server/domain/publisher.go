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

	p.pool.subscribersMutex.RLock()
	defer p.pool.subscribersMutex.RUnlock()
	for subscriber := range p.pool.subscribers {
		subscriber.receiveMessage(messageToSend)
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
