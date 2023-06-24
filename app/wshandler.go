package app

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/surasithaof/websocket-go/ws"
)

func wsHandler(hub *ws.Hub) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, err := hub.Connect(ctx.Writer, ctx.Request, uuid.NewString())
		if err != nil {
			log.Println("connect error", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		HealthEventType := "health"
		hub.SetupEventHandler(HealthEventType, func(event ws.Event, c ws.WebSocketClient) error {
			msg := map[string]any{
				"message": "pong",
				"time":    time.Now(),
			}

			msgJSON, err := json.Marshal(msg)
			if err != nil {
				log.Println("mashal json error", err)
				return err
			}

			pongMsg := ws.Event{
				Type:    HealthEventType,
				Payload: msgJSON,
			}

			c.SendMessage(pongMsg)
			log.Println(string(event.Payload))
			return nil
		})

		c.SendMessage("Salut!")
		ctx.Status(http.StatusOK)
	}
}
