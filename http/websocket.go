package http

import (
	"bytes"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
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

		messageToSend := bytes.Replace(message, []byte("?"), []byte("!"), -1)
		err = c.WriteMessage(mt, messageToSend)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

