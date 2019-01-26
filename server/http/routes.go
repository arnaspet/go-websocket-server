package http

import (
	"github.com/arnaspet/teso_task/server/domain"
	"github.com/sirupsen/logrus"
)

func (s *server) InitRoutes(logger *logrus.Logger) {
	connectionPool := domain.NewConnectionPool(logger)
	s.router.GET("/ws/publish", NewWebsocket(logger, connectionPool).PublisherWebsocketHandler)
}
