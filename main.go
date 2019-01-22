package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"teso_task/http"
)

func main() {
	port := flag.Uint("port", 3001, "Port on which http server is run")
	flag.Parse()

	s := http.NewServer(initLogger())
	s.Start(uint16(*port))
}

func initLogger() *logrus.Logger {
	logger := logrus.New()

	debug, ok := os.LookupEnv("DEBUG")

	if ok && strings.ToLower(debug) == "true" {
		logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}
