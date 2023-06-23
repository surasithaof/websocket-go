package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader
var allowedOrigins []string

type Config struct {
	AllowedOrigins  []string `envconfig:"WS_ALLOWED_ORIGINS" default:"*"`
	ReadBufferSize  int      `envconfig:"WS_READ_BUFFER_SIZE" default:"1024"`
	WriteBufferSize int      `envconfig:"WS_WRITE_BUFFER_SIZE" default:"1024"`
}

func Initialize(router *gin.RouterGroup, config Config) WebSocketClientManager {
	hub := NewHub()

	allowedOrigins = config.AllowedOrigins
	upgrader = websocket.Upgrader{
		ReadBufferSize:  config.ReadBufferSize,
		WriteBufferSize: config.WriteBufferSize,
		//Solving cross-domain problems
		CheckOrigin: checkOrigin,
	}

	return hub
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}

// func checkOrigin(r *http.Request) bool {
// 	return true
// }
