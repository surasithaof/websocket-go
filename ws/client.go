package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client
type Client struct {
	hub *Hub
	//Connection/Client ID
	id string
	//UserID
	userID string
	//Connected socket
	conn    *websocket.Conn
	message chan []byte
}

func newClient(conn *websocket.Conn, hub *Hub, userID string) WebSocketClient {
	c := Client{
		hub:     hub,
		id:      uuid.NewString(),
		userID:  userID,
		conn:    conn,
		message: make(chan []byte),
	}
	return &c
}

func (c *Client) GetClient() *Client {
	return c
}

func (c *Client) startReader() {
	defer func() {
		close(c.message)
		c.hub.removeClient(c)
	}()

	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Println(err)
		return
	}
	// Configure how to handle Pong responses
	c.conn.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read message client id:%s, error:%v", c.id, err)
			}
			break
		}

		var event Event
		err = json.Unmarshal(payload, &event)
		if err != nil {
			log.Println("marshal message error:", err)
			continue
		}

		err = c.hub.routeEvent(event, c)
		if err != nil {
			log.Println("handling message error:", err)
		}
	}
}

// pongHandler is used to handle PongMessages for the Client
func (c *Client) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	if logPingPongHealh {
		log.Println("pong", c.id)
	}
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}

func (c *Client) SendMessage(payload any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Printf("marshal payload error:%v", err)
		return err
	}

	c.message <- payloadJSON
	return nil
}

func (c *Client) startWriter() {
	// Create a ticker that triggers a ping at given interval
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.hub.removeClient(c)
	}()

	for {
		select {
		case msg, ok := <-c.message:
			// Ok will be false Incase the egress channel is closed
			if !ok {
				// Manager has closed this connection channel, so communicate that to frontend
				// Return to close the goroutine
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("send message error: ", err)
			}
			log.Println("message send to", c.id)
		case <-ticker.C:
			if logPingPongHealh {
				log.Println("ping", c.id)
			}
			// Send the Ping
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("writemsg: ", err)
				return // return to break this goroutine triggeing cleanup
			}
		}
	}
}
