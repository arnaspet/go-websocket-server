package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type server struct {
	logger *log.Logger
}

func New() *server {
	s := &server{
		log.New(os.Stdout, "http_server ", log.LstdFlags | log.Lshortfile),
	}

	return s
}

func (s *server) Start(port uint16) {
	s.logger.Printf("Starting server on %v\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)

	if err != nil {
		s.logger.Fatalf("Cannot run server: %v", err)
	}
}
