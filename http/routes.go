package http

func (s *server) routes() {
	s.router.GET("/ws", Websocket)
}
