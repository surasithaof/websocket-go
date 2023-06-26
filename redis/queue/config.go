package redisqueues

import (
	"crypto/tls"
	"time"
)

type ClientConfig struct {
	// // Network type to use, either tcp or unix.
	// // Default is tcp.
	// Network string

	// Redis server address in "host:port" format.
	// Addr string `envconfig:"REDIS_ADDRESS" required:"true"`

	// Username to authenticate the current connection when Redis ACLs are used.
	// See: https://redis.io/commands/auth.
	// Username string `envconfig:"REDIS_USERNAME"`

	// Password to authenticate the current connection.
	// See: https://redis.io/commands/auth.
	// Password string `envconfig:"REDIS_PASSWORD"`

	// Redis DB to select after connecting to a server.
	// See: https://redis.io/commands/select.
	// DB int `envconfig:"REDIS_DB" default:"0"`

	// ClientConfig
	URL string `envconfig:"REDIS_URL" required:"true"`

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration `envconfig:"REDIS_DIAL_TIMEOUT" default:"5s"`

	// Timeout for socket reads.
	// If timeout is reached, read commands will fail with a timeout error
	// instead of blocking.
	//
	// Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration `envconfig:"REDIS_READ_TIMEOUT" default:"3s"`

	// Timeout for socket writes.
	// If timeout is reached, write commands will fail with a timeout error
	// instead of blocking.
	//
	// Use value -1 for no timeout and 0 for default.
	// Default is ReadTimout.
	WriteTimeout time.Duration `envconfig:"REDIS_WRITE_TIMEOUT" default:"3s"`

	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int `envconfig:"REDIS_POOL_SIZE" default:"10"`

	// TLS Config used to connect to a server.
	// TLS will be negotiated only if this field is set.
	TLSConfig *tls.Config
}
