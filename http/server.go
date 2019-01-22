package http

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type server struct {
	logger *log.Logger
	router *httprouter.Router
}

func NewServer(logger *log.Logger) *server {
	s := &server{
		logger,
		httprouter.New(),
	}
	s.InitRoutes(logger)

	return s
}

func (s *server) Start(port uint16) {
	s.logger.Printf("Starting server on %v\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), s.router)

	if err != nil {
		s.logger.Fatalf("Cannot run server: %v", err)
	}
}
