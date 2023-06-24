package ws_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/surasithaof/websocket-go/ws"
)

const (
	TestPort         = "3002"
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
	done := make(chan bool)
	defer close(done)

	// setup gin router and websocket handler
	webSockets := ws.NewWebSocket(config)
	router := gin.Default()
	router.GET("/ws", wsHandler(webSockets.GetHub(), received, done))
	go router.Run(fmt.Sprintf(":%s", TestPort))
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// create websocket client to connect to server
	var client *websocket.Conn
	t.Run("connect to server", func(t *testing.T) {
		client, err = createClient()
		require.NoError(t, err)
		if err != nil {
			log.Fatal("dial error:", err)
			return
		}
	})
	defer func() {
		if client != nil {
			client.Close()
		}
	}()

	// test read message from the server
	t.Run("read message", func(t *testing.T) {
		go func() {
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
		}()
	})

	t.Run("send message", func(t *testing.T) {
		// test send message to the server
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
		case <-done:
			return
		case <-ticker.C:
			require.FailNow(t, "no signal received")
		}
	})
}

func createClient() (*websocket.Conn, error) {
	addr := flag.String("addr", fmt.Sprintf("localhost:%s", TestPort), "http service address")
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	return c, err
}

func wsHandler(hub *ws.Hub, msg chan ws.Event, done chan bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, err := hub.Connect(ctx.Writer, ctx.Request, uuid.NewString())
		if err != nil {
			log.Println("connect error", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		hub.SetupEventHandler(HealthCheckEvent, func(event ws.Event, c ws.WebSocketClient) error {
			msg <- event
			log.Println("got message", string(event.Payload))
			done <- true
			return nil
		})

		c.SendMessage("Salut!")
		ctx.Status(http.StatusOK)
	}
}
