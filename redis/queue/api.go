package redisqueues

import (
	"context"
	"time"
)

//go:generate mockgen -source=./api.go -destination=./mocks/queues.go -package=mocks "github.com/surasithaof/go-redis-async-queue" QueueWorker QueueProducer

// Worker processes jobs from a queue
type QueueWorker interface {
	Register(queueName string, handler func(c context.Context, payload []byte) error)
	Start()
	Shutdown()
}

// Producer sends jobs to a queue. User is responsible of encoding message to bytes
type QueueProducer interface {
	// Enqueue a job to run ASAP
	Enqueue(jobType string, payload []byte)
	// Schedule for later processing
	Schedule(jobType string, payload []byte, t time.Time)
	Shutdown() error
}
