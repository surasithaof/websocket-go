package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/surasithaof/websocket-go/adapters/httpserver"
	"github.com/surasithaof/websocket-go/ws"
)

func Start() error {
	_config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	config := _config
	router := httpserver.InitGin(config.HttpServer)

	httpserver.Run()
	rGroup := router.Group(config.HttpServer.Prefix)

	ws := ws.NewWebSocket(config.WebSocket)
	defer ws.Close()

	initialApp(rGroup, config, ws.GetHub())

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpserver.Shutdown(ctx)

	return nil
}

func initialApp(rGroup *gin.RouterGroup, config *Config, hub *ws.Hub) {
	rGroup.StaticFS("/client", http.Dir("./client"))
	rGroup.GET("/health", healthCheck)

	rGroup.GET("/ws", wsHandler(hub))

	// _ = db.Connect(config.Database)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func wsHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, err := hub.Connect(ctx.Writer, ctx.Request, uuid.NewString())
		if err != nil {
			log.Println("connect error", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		hub.SetupEventHandler("test", func(event ws.Event, c ws.WebSocketClient) error {
			c.SendMessage("got your message")
			log.Println(string(event.Payload))
			return nil
		})

		c.SendMessage("Salut!")
		ctx.Status(http.StatusOK)
	}
}
