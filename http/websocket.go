package http

import (
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"teso_task/domain"
)

var upgrader = websocket.Upgrader{}

func Websocket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer func() {
		if err := c.Close(); err != nil {
			log.Print("close:", err)
		}
	}()

	for {
		mt, message, err := c.ReadMessage()

		if mt == websocket.BinaryMessage {
			_ := c.Close()
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
