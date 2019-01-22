package http

import (
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"teso_task/domain"
	"time"
)

const (
	closeGracePeriod = 1000
)

type Websocket struct {
	logger *log.Logger
}

var upgrader = websocket.Upgrader{}

func NewWebsocket(logger *log.Logger) *Websocket {
	return &Websocket{
		logger,
	}
}

func (websocket *Websocket) WebsocketHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer closeWebsocketConnection(c)

	socketMessageLoop(c)
}

func socketMessageLoop(c *websocket.Conn) {
	for {
		mt, message, err := c.ReadMessage()

		if mt != websocket.TextMessage {
			closeWebsocketConnection(c)
			break
		}

		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		err = c.WriteMessage(mt, domain.ReplaceBytes(message))

		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

// gracefully close websocket connection
func closeWebsocketConnection(c *websocket.Conn) {
	if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.Print("close: ", err)
	}

	time.Sleep(closeGracePeriod)

	if err := c.Close(); err != nil {
		log.Print("close: ", err)
	}
}
