package http

import "github.com/sirupsen/logrus"

func (s *server) InitRoutes(logger *logrus.Logger) {
	s.router.GET("/ws", NewWebsocket(logger).WebsocketHandler)
}
