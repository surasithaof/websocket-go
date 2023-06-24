package app

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/surasithaof/websocket-go/adapters/db"
	"github.com/surasithaof/websocket-go/adapters/httpserver"
	"github.com/surasithaof/websocket-go/ws"
)

type Config struct {
	HttpServer httpserver.Config
	Database   db.Config
	WebSocket  ws.Config
}

func LoadConfig() (*Config, error) {
	AppConfig := &Config{}
	err := envconfig.Process("", AppConfig)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	return AppConfig, nil
}
