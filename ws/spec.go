package ws

import (
	"net/http"
)

//go:generate mockgen -source=./spec.go -destination=./mocks/websocket.go -package=mocks "surasithaof/websocket-go" WebSocketClient WebSocketClientManager

type WebSocketClient interface {
	GetClient() *Client
	SendMessage(payload any) error
}

type WebSocketClientManager interface {
	GetHub() *Hub
	Connect(w http.ResponseWriter, r *http.Request, userID string) (WebSocketClient, error)
	SetupEventHandler(eventType string, handlerFunc EventHandler)
	GetClient(ID string) (bool, WebSocketClient)
	GetClientsByUserID(userID string) []WebSocketClient
	Broadcast(payload any) error
	Close()
}
