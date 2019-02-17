package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	addr  = flag.String("addr", "localhost:3001", "http service address")
	path  = flag.String("path", "/ws/subscribe", "http service path for subscribing")
	debug = flag.Bool("debug", false, "Should output be debugged")
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	logger := initLogger(*debug)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: *path}
	logger.Infof("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Fatal("dial: ", err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			logger.Fatal("close: ", err)
		}
	}()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logger.Error("read: ", err)
				return
			}
			logger.Debugf("Receive: %s", message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			logger.Info("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Error("write close: ", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func initLogger(debug bool) *logrus.Logger {
	logger := logrus.New()

	if debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}
