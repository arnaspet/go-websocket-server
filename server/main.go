package main

import (
	"flag"
	"github.com/arnaspet/teso_task/server/http"
	"github.com/sirupsen/logrus"
)

var (
	port = flag.Uint("port", 3001, "Port on which http server is run")
	debug = flag.Bool("debug", false, "Should output be debugged")
)

func main() {
	flag.Parse()

	s := http.NewServer(initLogger(*debug))
	s.Start(uint16(*port))
}

func initLogger(debug bool) *logrus.Logger {
	logger := logrus.New()

	if debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}
