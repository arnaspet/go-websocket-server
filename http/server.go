package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/julienschmidt/httprouter"
)

type server struct {
	logger *log.Logger
	router *httprouter.Router
}

func New() *server {
	s := &server{
		log.New(os.Stdout, "http_server ", log.LstdFlags | log.Lshortfile),
		httprouter.New(),
	}
	s.routes()

	return s
}

func (s *server) Start(port uint16) {
	s.logger.Printf("Starting server on %v\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), s.router)

	if err != nil {
		s.logger.Fatalf("Cannot run server: %v", err)
	}
}
