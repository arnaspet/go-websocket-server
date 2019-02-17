package domain

import (
	"github.com/gorilla/websocket"
)

type Publisher struct {
	pool *ConnectionPool
	conn *Connection
}

func (p *Publisher) getConnection() *websocket.Conn {
	return p.conn.websocketConn
}

func (p *Publisher) getId() uint {
	return p.conn.id
}

func (p *Publisher) receiveMessage(message *Message) {
	p.conn.receiveMessage(message)
}

func(p *Publisher) broadcastMessage(message *Message) {
	messageToSend := &Message{
		ReplaceBytes(message.content),
		websocket.TextMessage,
	}

	p.pool.subscribersMutex.RLock()
	defer p.pool.subscribersMutex.RUnlock()
	for subscriber := range p.pool.subscribers {
		subscriber.receiveMessage(messageToSend)
	}
}

func (p *Publisher) initMessageHandler() {
	for {
		mt, message, err := p.getConnection().ReadMessage()

		if err != nil {
			p.conn.logger.Error("read: ", err)
			break
		}

		if mt != websocket.TextMessage {
			p.pool.closeWebsocketConnection(p)
			break
		}

		p.conn.logger.Debugf("Message from publisher #%d received: %s", p.conn.id, message)
		p.broadcastMessage(&Message{
			message,
			websocket.TextMessage,
		})
	}
}
