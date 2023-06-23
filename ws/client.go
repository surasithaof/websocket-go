package ws

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client
type Client struct {
	Hub *Hub
	//Connection/Client ID
	ID string
	//UserID
	UserID string
	//Connected socket
	Conn *websocket.Conn
}

func NewClient(conn *websocket.Conn, hub *Hub, userID string) WebSocketClient {
	c := Client{
		Hub:    hub,
		ID:     uuid.NewString(),
		UserID: userID,
		Conn:   conn,
	}
	return &c
}

func (c *Client) GetClient() *Client {
	return c
}

func (c *Client) ReadMessage(handlerFunc func(payload []byte) error) {
	defer c.Hub.removeClient(c)
	for {
		_, payload, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("read message client id:%s, error:%v", c.ID, err)
			return
		}
		handlerFunc(payload)
	}
}

func (c *Client) SendMessage(payload any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Printf("marshal payload error:%v", err)
		return err
	}
	err = c.Conn.WriteMessage(websocket.TextMessage, payloadJSON)
	if err != nil {
		log.Printf("write message client id:%s, error:%v", c.ID, err)
		return err
	}
	log.Println("message send to", c.ID)
	return nil
}
