package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader

type Config struct {
	AllowedOrigins  []string `envconfig:"WS_ALLOWED_ORIGINS" default:"*"`
	ReadBufferSize  int      `envconfig:"WS_READ_BUFFER_SIZE" default:"1024"`
	WriteBufferSize int      `envconfig:"WS_WRITE_BUFFER_SIZE" default:"1024"`
}

func NewWebSocket(config Config) WebSocketClientManager {
	hub := newHub()

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
