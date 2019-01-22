package http

import "log"

func (s *server) InitRoutes(logger *log.Logger) {
	s.router.GET("/ws", NewWebsocket(logger).WebsocketHandler)
}
