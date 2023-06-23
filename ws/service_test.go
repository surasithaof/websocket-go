package ws_test

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func HandlerTest(c *gin.Context) {
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	for {
		//Read Message from client
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("received message:", string(message))

		//If client message is ping will return pong
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//Response message to client
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func TestConnectWebSocket(t *testing.T) {
	r := gin.Default()
	r.GET("/ws", HandlerTest)
	go r.Run(":3002")

	addr := flag.String("addr", "localhost:3002", "http service address")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	require.NoError(t, err)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	log.Printf("connect to %s success\n", u.String())
	defer c.Close()

	done := make(chan struct{})
	msg := make(chan string)
	defer close(msg)

	go func() {
		defer close(done)
		for {
			mt, message, err := c.ReadMessage()
			require.NoError(t, err)
			if err != nil {
				log.Fatal("send message error:", err)
				return
			}
			log.Println("message type:", mt)
			log.Println("message:", string(message))

			msg <- string(message)
		}
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte("ping"))
	if err != nil {
		log.Println("write:", err)
		return
	}

	timeout := 1 * time.Second
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	select {
	case receivedMessage := <-msg:
		assert.Equal(t, "pong", receivedMessage)
		return
	case <-done:
		return
	case <-ticker.C:
		require.FailNow(t, "no signal received")
	}
}
