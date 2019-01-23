package http

import (
	"github.com/arnaspet/teso_task/domain"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	closeGracePeriod = 1000
)

type Websocket struct {
	logger *logrus.Logger
}

var upgrader = websocket.Upgrader{}

func NewWebsocket(logger *logrus.Logger) *Websocket {
	return &Websocket{
		logger,
	}
}

func (ws *Websocket) WebsocketHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.logger.Error("upgrade: ", err)
		return
	}

	defer ws.closeWebsocketConnection(c)

	ws.socketMessageLoop(c)
}

func (ws *Websocket) socketMessageLoop(c *websocket.Conn) {
	for {
		mt, message, err := c.ReadMessage()

		if mt != websocket.TextMessage {
			ws.closeWebsocketConnection(c)
			break
		}

		if err != nil {
			ws.logger.Error("read: ", err)
			break
		}
		ws.logger.Debugf("recv: %s", message)

		err = c.WriteMessage(mt, domain.ReplaceBytes(message))

		if err != nil {
			ws.logger.Error("write: ", err)
			break
		}
	}
}

func (ws *Websocket) closeWebsocketConnection(c *websocket.Conn) {
	ws.logger.Debug("Gracefully closing websocket connection")
	if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		ws.logger.Warn("close: ", err)
	}

	time.Sleep(closeGracePeriod)

	if err := c.Close(); err != nil {
		ws.logger.Warn("close: ", err)
	}
}
