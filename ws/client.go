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
	Conn    *websocket.Conn
	message chan []byte
}

func newClient(conn *websocket.Conn, hub *Hub, userID string) WebSocketClient {
	c := Client{
		Hub:     hub,
		ID:      uuid.NewString(),
		UserID:  userID,
		Conn:    conn,
		message: make(chan []byte),
	}
	return &c
}

func (c *Client) GetClient() *Client {
	return c
}

func (c *Client) ReadMessages(handlerFunc func(payload []byte) error) {
	defer func() {
		c.Hub.removeClient(c)
	}()

	for {
		_, payload, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read message client id:%s, error:%v", c.ID, err)
			}
			break
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

	c.message <- payloadJSON
	return nil
}

func (c *Client) startWriter() {
	defer func() {
		c.Hub.removeClient(c)
	}()

	for msg := range c.message {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("send message error: ", err)
		}
		log.Println("message send to", c.ID)
	}
}
