package http

import (
	"github.com/arnaspet/teso_task/server/domain"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)


type Websocket struct {
	logger *logrus.Logger
	pool *domain.ConnectionPool
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

func (ws *Websocket) socketMessageLoop(c *websocket.Conn) {
	//for {
	//	mt, message, err := c.ReadMessage()
	//
	//	if mt != websocket.TextMessage {
	//		ws.closeWebsocketConnection(c)
	//		break
	//	}
	//
	//	if err != nil {
	//		ws.logger.Error("read: ", err)
	//		break
	//	}
	//	ws.logger.Debugf("recv: %s", message)
	//
	//	err = c.WriteMessage(mt, domain.ReplaceBytes(message))
	//
	//	if err != nil {
	//		ws.logger.Error("write: ", err)
	//		break
	//	}
	//}
}
