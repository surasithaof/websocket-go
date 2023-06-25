package ws

import (
	"encoding/json"
	"log"

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
		c.hub.removeClient(c)
	}()

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
	defer func() {
		c.hub.removeClient(c)
	}()

	for msg := range c.message {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("send message error: ", err)
		}
		log.Println("message send to", c.id)
	}
}
