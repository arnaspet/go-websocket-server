package domain

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

const (
	noSubscribersCode        = 1
	subscribersListeningCode = 2
)

type Publisher struct {
	pool *ConnectionPool
	conn *Connection
}

type publisherNotification struct {
	Code    int
	Message []byte
}

func (p *Publisher) getConnection() *websocket.Conn {
	return p.conn.websocketConn
}

func (p *Publisher) getId() uint64 {
	return p.conn.id
}

func (p *Publisher) receiveMessage(message *Message) {
	p.conn.receiveMessage(message)
}

func (p *Publisher) broadcastMessage(message *Message) {
	messageToSend := &Message{
		ReplaceBytes(message.Content),
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
		_, message, err := p.getConnection().ReadMessage()

		if err != nil {
			p.conn.logger.Error("read: ", err)
			break
		}

		p.conn.logger.Debugf("Message from publisher #%d received: %s", p.conn.id, message)
		p.broadcastMessage(&Message{
			message,
			websocket.TextMessage,
		})
	}
}

func (p *Publisher) Notify(code int) {
	content, err := json.Marshal(createPublisherNotification(code))

	if err != nil {
		p.conn.logger.Error("json: ", err)
	}

	p.receiveMessage(&Message{
		Content: content,
		MsgType: websocket.TextMessage,
	})
}

func createPublisherNotification(code int) *publisherNotification {
	var notification *publisherNotification

	switch code {
	case noSubscribersCode:
		notification = &publisherNotification{noSubscribersCode, []byte("No subscribers are currently listening")}
	case subscribersListeningCode:
		notification = &publisherNotification{subscribersListeningCode, []byte("Subscribers started listening")}
	}

	return notification
}
