package app

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/surasithaof/websocket-go/adapters/db"
	"github.com/surasithaof/websocket-go/adapters/httpserver"
)

type Config struct {
	HttpServer httpserver.Config
	Database   db.Config
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
