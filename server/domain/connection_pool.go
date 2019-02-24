package domain

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

const (
	closeGracePeriod = 1000
)

type ConnectionPool struct {
	logger *logrus.Logger
	seq    uint64

	publishers      map[*Publisher]struct{}
	publishersMutex *sync.RWMutex

	subscribers      map[*Subscriber]struct{}
	subscribersMutex *sync.RWMutex
}

func NewConnectionPool(logger *logrus.Logger) *ConnectionPool {
	return &ConnectionPool{
		logger:           logger,
		publishers:       make(map[*Publisher]struct{}),
		publishersMutex:  &sync.RWMutex{},
		subscribers:      make(map[*Subscriber]struct{}),
		subscribersMutex: &sync.RWMutex{},
	}
}

func (cp *ConnectionPool) InitPublisher(conn *websocket.Conn) {
	id := cp.generateId()
	publisher := &Publisher{
		cp,
		NewConnection(conn, cp.logger, id, make(chan *Message)),
	}
	defer cp.closeWebsocketConnection(publisher)
	cp.registerPublisherToPool(publisher)

	publisher.initMessageHandler()
}

func (cp *ConnectionPool) registerPublisherToPool(publisher *Publisher) {
	cp.publishersMutex.Lock()
	defer cp.publishersMutex.Unlock()
	cp.publishers[publisher] = struct{}{}

	if len(cp.subscribers) > 0 {
		cp.notifyPublishers(subscribersListeningCode)
	}

	cp.logger.Debugf("#%d Publisher registered to pool", publisher.conn.id)
}

func (cp *ConnectionPool) InitSubscriber(conn *websocket.Conn) {
	id := cp.generateId()
	subscriber := &Subscriber{
		cp,
		NewConnection(conn, cp.logger, id, make(chan *Message)),
	}
	defer cp.closeWebsocketConnection(subscriber)
	cp.registerSubscriberToPool(subscriber)

	subscriber.initMessageHandler()
}

func (cp *ConnectionPool) registerSubscriberToPool(subscriber *Subscriber) {
	cp.subscribersMutex.Lock()
	defer cp.subscribersMutex.Unlock()
	cp.subscribers[subscriber] = struct{}{}

	if len(cp.subscribers) == 1 {
		cp.notifyPublishers(subscribersListeningCode)
	}

	cp.logger.Debugf("#%d Subscriber registered to pool", subscriber.conn.id)
}

func (cp *ConnectionPool) notifyPublishers(code int) {
	for publisher := range cp.publishers {
		publisher.Notify(code)
	}
}

func (cp *ConnectionPool) closeWebsocketConnection(ch ConnectionHolder) {
	cp.unregisterFromPool(ch)
	cp.logger.Debugf("Gracefully closing websocket connection %v", ch.getId())
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

		if len(cp.subscribers) == 0 {
			cp.notifyPublishers(noSubscribersCode)
		}
	}
}

func (cp *ConnectionPool) generateId() uint64 {
	atomic.AddUint64(&cp.seq, 1)

	return cp.seq
}
