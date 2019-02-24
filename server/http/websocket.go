package http

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/arnaspet/teso_task/server/domain"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type Websocket struct {
	logger *logrus.Logger
	pool   *domain.ConnectionPool
}

var upgrader = websocket.Upgrader{}

func NewWebsocket(logger *logrus.Logger, pool *domain.ConnectionPool) *Websocket {
	return &Websocket{
		logger,
		pool,
	}
}

func (ws *Websocket) PublisherWebsocketHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.logger.Error("upgrade: ", err)
		return
	}

	ws.pool.InitPublisher(c)
}

func (ws *Websocket) SubscriberWebsocketHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.logger.Error("upgrade: ", err)
		return
	}

	ws.pool.InitSubscriber(c)
}
