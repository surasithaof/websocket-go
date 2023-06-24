package ws

import (
	"log"
	"net/http"
	"sync"
)

// Hub for manage clients
type Hub struct {
	clients map[string]*Client

	// to handle incoming message with specific handler
	handlers map[string]EventHandler

	// Using a syncMutex here to be able to lcok state before editing clients
	// Could also use Channels to block
	sync.RWMutex
}

func newHub() WebSocketClientManager {
	hub := Hub{
		clients:  make(map[string]*Client),
		handlers: make(map[string]EventHandler),
	}
	return &hub
}

func (h *Hub) SetupEventHandler(eventType string, handlerFunc EventHandler) {
	h.handlers[eventType] = handlerFunc
}

func (h *Hub) GetHub() *Hub {
	return h
}

func (h *Hub) Connect(w http.ResponseWriter, r *http.Request, userID string) (WebSocketClient, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket connect error :", err)
		return nil, err
	}
	client := newClient(conn, h, userID)
	c := client.GetClient()
	h.addClient(c.GetClient())

	go c.startWriter()

	go c.startReader()

	return client, nil
}

func (h *Hub) GetClient(ID string) (bool, *Client) {
	c, found := h.clients[ID]
	if !found {
		return false, nil
	}
	return true, c
}

func (h *Hub) GetClientsByUserID(userID string) []*Client {
	userClients := []*Client{}

	for _, c := range h.clients {
		if c.UserID == userID {
			userClients = append(userClients, c)
		}
	}

	return userClients
}

func (h *Hub) Broadcast(payload any) error {
	for _, c := range h.clients {
		err := c.SendMessage(payload)
		if err != nil {
			log.Printf("broadcast message to clientID:%s, error:%v", c.ID, err)
			return err
		}
	}
	return nil
}

func (h *Hub) addClient(client *Client) *Client {
	h.Lock()
	defer h.Unlock()

	h.clients[client.ID] = client
	return client
}

func (h *Hub) removeClient(client *Client) {
	h.Lock()
	defer h.Unlock()

	err := client.Conn.Close()
	if err != nil {
		log.Printf("close message client id:%s, error:%v", client.ID, err)
	}
	delete(h.clients, client.ID)
	log.Println("closed client id: ", client.ID)
}

func (h *Hub) Close() {
	for _, c := range h.clients {
		h.removeClient(c)
	}
}

func (h *Hub) routeEvent(event Event, c *Client) error {
	handler, ok := h.handlers[event.Type]
	if !ok {
		return ErrEventNotSupported
	}

	err := handler(event, c)
	if err != nil {
		return err
	}
	return nil
}
