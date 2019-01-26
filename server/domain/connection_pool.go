package domain

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ConnectionPool struct {
	logger *logrus.Logger
	publishers []*Publisher
	seq uint
}

func NewConnectionPool(logger *logrus.Logger) *ConnectionPool {
	return &ConnectionPool{
		logger:     logger,
		publishers: make([]*Publisher, 0),
		seq:        0,
	}
}

func(cp *ConnectionPool) RegisterConnectionAsPublisher(conn *websocket.Conn) {
	id := cp.generateId()
	publisher := &Publisher{
		conn,
		id,
	}

	cp.publishers = append(cp.publishers, publisher)
	cp.logger.Debugf("#%d Publisher registered to pool", id)
}

func(cp *ConnectionPool) generateId() uint {
	cp.seq += 1

	return cp.seq
}
