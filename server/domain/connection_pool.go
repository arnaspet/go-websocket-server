package domain

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	closeGracePeriod = 1000
)

type ConnectionPool struct {
	logger *logrus.Logger
	seq uint

	publishers map[*Publisher]struct{}
	publishersMutex *sync.RWMutex

	subscribers map[*Subscriber]struct{}
	subscribersMutex *sync.RWMutex
}

func NewConnectionPool(logger *logrus.Logger) *ConnectionPool {
	return &ConnectionPool{
		logger:     logger,
		publishers: make(map[*Publisher]struct{}),
		publishersMutex: &sync.RWMutex{},
		subscribers: make(map[*Subscriber]struct{}),
		subscribersMutex: &sync.RWMutex{},
	}
}

func(cp *ConnectionPool) InitPublisher(conn *websocket.Conn) {
	id := cp.generateId()
	publisher := &Publisher{
		cp,
		&Connection{
			conn,
			cp.logger,
			id,
			make(chan *Message),
		},
	}
	defer cp.closeWebsocketConnection(publisher)

	cp.publishersMutex.Lock()
	cp.publishers[publisher] = struct{}{}
	cp.publishersMutex.Unlock()
	cp.logger.Debugf("#%d Publisher registered to pool", id)

	publisher.initMessageHandler()
}


func(cp *ConnectionPool) InitSubscriber(conn *websocket.Conn) {
	id := cp.generateId()
	subscriber := &Subscriber{
		cp,
		&Connection{
			conn,
			cp.logger,
			id,
			make(chan *Message),
		},
	}
	defer cp.closeWebsocketConnection(subscriber)

	cp.subscribersMutex.Lock()
	cp.subscribers[subscriber] = struct{}{}
	cp.subscribersMutex.Unlock()
	cp.logger.Debugf("#%d Subscriber registered to pool", id)

	subscriber.conn.initMessageSender()
	subscriber.initMessageHandler()
}

func (cp *ConnectionPool) closeWebsocketConnection(ch ConnectionHolder) {
	cp.unregisterFromPool(ch)
	cp.logger.Debug("Gracefully closing websocket connection")
	closeMessage := &Message{
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		websocket.CloseMessage,
	}

	ch.receiveMessage(closeMessage)
	time.Sleep(closeGracePeriod)

	if err := ch.getConnection().Close(); err != nil {
		cp.logger.Warn("close: ", err)
	}
}

func (cp *ConnectionPool) unregisterFromPool(ch ConnectionHolder) {
	switch ch.(type) {
	case *Publisher:
		cp.publishersMutex.Lock()
		defer cp.publishersMutex.Unlock()
		delete(cp.publishers, ch.(*Publisher))

	case *Subscriber:
		cp.subscribersMutex.Lock()
		defer cp.subscribersMutex.Unlock()
		delete(cp.subscribers, ch.(*Subscriber))
	}
}

func(cp *ConnectionPool) generateId() uint {
	cp.seq += 1

	return cp.seq
}
