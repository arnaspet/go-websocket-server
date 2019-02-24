package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/gorilla/websocket"
)

var (
	addr  = flag.String("addr", "localhost:3001", "http service address")
	path  = flag.String("path", "/ws/publish", "http service path for publishing")
	debug = flag.Bool("debug", false, "Should output be debugged")
)

func main() {
	flag.Parse()
	logger := initLogger(*debug)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c := initWebsocketConnection(logger)

	done := make(chan struct{})
	ticker := time.NewTicker(time.Second)

	publisher := NewPublisher(logger, c, done, interrupt, ticker)
	publisher.Start()
}

func initWebsocketConnection(logger *logrus.Logger) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: *addr, Path: *path}

	logger.Infof("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	return c
}

func initLogger(debug bool) *logrus.Logger {
	logger := logrus.New()

	if debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}
