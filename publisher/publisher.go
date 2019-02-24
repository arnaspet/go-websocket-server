package main

import (
	"encoding/json"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/icrowley/fake"
)

const (
	stateSending = 1
	stateWaiting = 2

	noSubscribersCode        = 1
	subscribersListeningCode = 2
)

type publisher struct {
	logger *logrus.Logger
	conn   *websocket.Conn

	done      chan struct{}
	interrupt chan os.Signal
	ticker    *time.Ticker

	state int32
}

func NewPublisher(logger *logrus.Logger, conn *websocket.Conn, done chan struct{}, interrupt chan os.Signal, ticker *time.Ticker) *publisher {
	publisher := &publisher{
		logger,
		conn,
		done,
		interrupt,
		ticker,
		stateWaiting,
	}

	return publisher
}

func (p *publisher) Start() {
	defer func() {
		err := p.conn.Close()

		if err != nil {
			p.logger.Error("close error: ", err)
		}
	}()
	defer p.ticker.Stop()

	p.initReceiveLoop()
	p.initMainLoop()
}

func (p *publisher) initMainLoop() {
	for {
		select {
		case <-p.done:
			return
		case <-p.ticker.C:
			if atomic.LoadInt32(&p.state) == stateWaiting {
				continue
			}

			message := []byte(fake.Color())
			p.logger.Infof("Sending message: %s", message)
			err := p.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				p.logger.Error("write:", err)
				return
			}
		case <-p.interrupt:
			p.logger.Infoln("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := p.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				p.logger.Error("write close:", err)
				return
			}
			select {
			case <-p.done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (p *publisher) initReceiveLoop() {
	go func() {
		defer close(p.done)
		for {
			_, message, err := p.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			var f interface{}
			err = json.Unmarshal(message, &f)
			decodedMessage := f.(map[string]interface{})

			if err != nil {
				log.Println("json:", err)
			}

			switch int(decodedMessage["Code"].(float64)) {
			case noSubscribersCode:
				log.Println("Got do not send command!")
				atomic.StoreInt32(&p.state, stateWaiting)

			case subscribersListeningCode:
				log.Println("Got a send command!")
				atomic.StoreInt32(&p.state, stateSending)
			}
		}
	}()
}
