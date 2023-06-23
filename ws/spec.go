package ws

import (
	"net/http"
)

type WebSocketClient interface {
	GetClient() *Client
	ReadMessage(handlerFunc HandlerFunc)
	SendMessage(payload any) error
}

type HandlerFunc = func(payload []byte) error

type WebSocketClientManager interface {
	GetHub() *Hub
	Connect(w http.ResponseWriter, r *http.Request, userID string) (WebSocketClient, error)
	GetClient(ID string) (bool, *Client)
	GetClientsByUserID(userID string) []*Client
	Broadcast(payload any) error
	Close()
}
