package ws

import (
	"encoding/json"
	"errors"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c WebSocketClient) error

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)
