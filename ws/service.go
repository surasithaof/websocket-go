package ws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader

type Config struct {
	AllowedOrigins   []string `envconfig:"WS_ALLOWED_ORIGINS" default:"*"`
	ReadBufferSize   int      `envconfig:"WS_READ_BUFFER_SIZE" default:"1024"`
	WriteBufferSize  int      `envconfig:"WS_WRITE_BUFFER_SIZE" default:"1024"`
	PongWaitSec      int      `envconfig:"WS_PONG_WAIT_SECS" default:"10"`
	LogPingPongHealh bool     `envconfig:"WS_LOG_PONG_PONG_HEALTH" default:"false"`
}

var (
	logPingPongHealh = false
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

func NewWebSocket(config Config) WebSocketClientManager {
	hub := newHub()

	pongWait = time.Second * time.Duration(config.PongWaitSec)
	pingInterval = (pongWait * 9) / 10
	logPingPongHealh = config.LogPingPongHealh

	upgrader = websocket.Upgrader{
		ReadBufferSize:  config.ReadBufferSize,
		WriteBufferSize: config.WriteBufferSize,
		//Solving cross-domain problems
		CheckOrigin: checkOrigun(config.AllowedOrigins),
	}

	return hub
}

func checkOrigun(allowedOrigins []string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		for _, allowed := range allowedOrigins {
			if allowed == "*" {
				return true
			}

			if origin == allowed {
				return true
			}
		}
		return false
	}
}
