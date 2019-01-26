package http

import (
	"github.com/arnaspet/teso_task/server/domain"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestWebsocketUpgradesToVersionAtleast13(t *testing.T) {
	req, _ := http.NewRequest("GET", "/ws", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, _ := test.NewNullLogger()
		pool := domain.NewConnectionPool(logger)

		NewWebsocket(logger, pool).PublisherWebsocketHandler(w, r, httprouter.Params{})
	})
	handler.ServeHTTP(rr, req)
	socketVersion, err := strconv.Atoi(rr.Header()["Sec-Websocket-Version"][0])
	assert.NoError(t, err)

	var comparison assert.Comparison = func() (success bool) {
		return socketVersion >= 13
	}

	assert.Condition(t, comparison)
}

func TestWebsocketClosesConnectionOnBinaryData(t *testing.T) {
	server := createHttpServer()
	defer server.Close()

	ws := getWebsocketConnection(t, server)

	ws.SetCloseHandler(func(code int, text string) error {
		assert.Equal(t, websocket.CloseNormalClosure, code)
		return nil
	})

	err := ws.WriteMessage(websocket.BinaryMessage, []byte("hello"))
	assert.NoError(t, err)

	_, _, err = ws.ReadMessage()
	assert.Error(t, err)
}

func TestWebsocketClosesAfterClientInitiatesClose(t *testing.T) {
	server := createHttpServer()
	defer server.Close()

	ws := getWebsocketConnection(t, server)
	assert.NoError(t, ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")))
}

func createHttpServer() *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, _ := test.NewNullLogger()
		pool := domain.NewConnectionPool(logger)

		NewWebsocket(logger, pool).PublisherWebsocketHandler(w, r, httprouter.Params{})
	}))

	return s
}

func getWebsocketConnection(t *testing.T, server *httptest.Server) *websocket.Conn {
	// Convert http://127.0.0.1 to ws://127.0.0.1
	u := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	return ws
}
