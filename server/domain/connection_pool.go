package domain

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	closeGracePeriod = 1000
)

type ConnectionHolder interface {
	getConnection() *websocket.Conn
	getId() uint
}

type ConnectionPool struct {
	logger *logrus.Logger
	publishers []*Publisher
	subscribers []*Subscriber
	seq uint
}

func NewConnectionPool(logger *logrus.Logger) *ConnectionPool {
	return &ConnectionPool{
		logger:     logger,
		publishers: make([]*Publisher, 0),
		seq:        0,
	}
}

func(cp *ConnectionPool) InitPublisher(conn *websocket.Conn) {
	id := cp.generateId()
	publisher := &Publisher{
		cp,
		conn,
		cp.logger,
		id,
	}
	defer cp.closeWebsocketConnection(publisher)

	cp.publishers = append(cp.publishers, publisher)
	cp.logger.Debugf("#%d Publisher registered to pool", id)

	publisher.initMessageHandler()
}


func(cp *ConnectionPool) InitSubscriber(conn *websocket.Conn) {
	id := cp.generateId()
	subscriber := &Subscriber{
		cp,
		conn,
		cp.logger,
		id,
	}
	defer cp.closeWebsocketConnection(subscriber)
	cp.subscribers = append(cp.subscribers, subscriber)
	cp.logger.Debugf("#%d Subscriber registered to pool", id)

	subscriber.initMessageHandler()
}

func (cp *ConnectionPool) closeWebsocketConnection(ch ConnectionHolder) {
	cp.logger.Debug("Gracefully closing websocket connection")
	if err := ch.getConnection().WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		cp.logger.Warn("close: ", err)
	}

	time.Sleep(closeGracePeriod)

	if err := ch.getConnection().Close(); err != nil {
		cp.logger.Warn("close: ", err)
	}
	cp.unregisterFromPool(ch)
}

func (cp *ConnectionPool) unregisterFromPool(ch ConnectionHolder) {
	for i := range cp.subscribers {
		if cp.subscribers[i].getId() == ch.getId() {
			cp.subscribers[i] = cp.subscribers[len(cp.subscribers)-1]
			cp.subscribers = cp.subscribers[:len(cp.subscribers)-1]
			return
		}
	}

	for i := range cp.publishers {
		if cp.publishers[i].getId() == ch.getId() {
			cp.publishers[i] = cp.publishers[len(cp.publishers)-1]
			cp.publishers = cp.publishers[:len(cp.publishers)-1]
			return
		}
	}
}

func(cp *ConnectionPool) generateId() uint {
	cp.seq += 1

	return cp.seq
}
