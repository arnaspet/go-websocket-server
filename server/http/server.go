package http

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

type server struct {
	logger *logrus.Logger
	router *httprouter.Router
}

func NewServer(logger *logrus.Logger) *server {
	s := &server{
		logger,
		httprouter.New(),
	}
	s.InitRoutes(logger)

	return s
}

func (s *server) Start(port uint16) {
	s.logger.Infof("Starting server on %v", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), s.router)

	if err != nil {
		s.logger.Fatalf("Cannot run server: %v", err)
	}
}
