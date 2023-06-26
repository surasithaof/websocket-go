package ws_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/surasithaof/websocket-go/ws"
)

const (
	TestPort         = 3002
	HealthCheckEvent = "health"
)

func TestWebSocket(t *testing.T) {
	config := ws.Config{}
	err := envconfig.Process("", &config)
	if err != nil {
		log.Println("load config error", err)
		t.Skip()
	}

	received := make(chan ws.Event)
	defer close(received)

	// setup http server and websocket handler
	webSockets := ws.NewWebSocket(config)
	defer webSockets.Close()

	s := httptest.NewServer(wsHandler(webSockets.GetHub(), received))
	defer s.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	url := "ws" + strings.TrimPrefix(s.URL, "http")
	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	require.NoError(t, err)
	if err != nil {
		log.Fatal("dial error:", err)
		return
	}
	defer client.Close()

	// test read message from the server
	t.Run("read message", func(t *testing.T) {

		_, payload, err := client.ReadMessage()
		require.NoError(t, err)
		if err != nil {
			log.Println("read message error", err)
			return
		}
		log.Println("client got message", string(payload))

		expected, err := json.Marshal("Salut!")
		require.NoError(t, err)
		if err != nil {
			log.Println("mashal expect message error", err)
			return
		}
		assert.Equal(t, expected, payload)

	})

	// test send message to the server
	t.Run("send message", func(t *testing.T) {

		healPayload := map[string]any{
			"message": "ping",
			"send":    time.Now(),
		}
		healPayloadJSON, err := json.Marshal(&healPayload)
		require.NoError(t, err)
		if err != nil {
			log.Println("mashal event payload error", err)
			return
		}

		healthEvent := ws.Event{
			Type:    HealthCheckEvent,
			Payload: healPayloadJSON,
		}
		testMsg, err := json.Marshal(&healthEvent)
		require.NoError(t, err)
		if err != nil {
			log.Println("mashal event message error", err)
			return
		}

		err = client.WriteMessage(websocket.TextMessage, testMsg)
		if err != nil {
			log.Println("send message error", err)
			return
		}

		timeout := 5 * time.Second
		ticker := time.NewTicker(timeout)
		defer ticker.Stop()

		select {
		case receivedMessage := <-received:
			assert.Equal(t, HealthCheckEvent, receivedMessage.Type)
			return
		case <-ticker.C:
			require.FailNow(t, "no signal received")
			return
		}

	})
}

func wsHandler(hub *ws.Hub, msg chan ws.Event) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := hub.Connect(w, r, uuid.NewString())
		if err != nil {
			log.Println("connect error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		hub.SetupEventHandler(HealthCheckEvent, func(event ws.Event, c ws.WebSocketClient) error {
			msg <- event
			log.Println("got message", string(event.Payload))
			return nil
		})

		c.SendMessage("Salut!")
		w.WriteHeader(http.StatusOK)
	}
}
