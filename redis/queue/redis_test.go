package redisqueues_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	redisqueues "github.com/surasithaof/websocket-go/redis/queue"
)

const (
	redisUrl = "redis://@localhost:6379"
)

func TestFull(t *testing.T) {
	if _, found := os.LookupEnv("E2E"); !found {
		t.Skip()
	}
	payload, _ := json.Marshal("message")
	queueName := "foo"

	clientConfig := redisqueues.ClientConfig{URL: redisUrl}
	producer := redisqueues.NewProducer(clientConfig)
	defer producer.Shutdown()
	worker := redisqueues.NewWorker(redisqueues.WorkerConfig{
		URL: redisUrl,
	})
	defer worker.Shutdown()

	callChan := make(chan any)
	worker.Register(queueName, func(c context.Context, received []byte) error {
		assert.Equal(t, payload, received)
		callChan <- true
		return nil
	})

	worker.Start()

	producer.Enqueue(queueName, payload)

	assertSignal(t, 3*time.Second, callChan, "handler not called")
}

func TestSchedule(t *testing.T) {
	if _, found := os.LookupEnv("E2E"); !found {
		t.Skip()
	}
	queueName := "foo"
	payload, _ := json.Marshal("message")

	clientConfig := redisqueues.ClientConfig{URL: redisUrl}
	producer := redisqueues.NewProducer(clientConfig)
	defer producer.Shutdown()
	worker := redisqueues.NewWorker(redisqueues.WorkerConfig{
		URL: redisUrl,
	})
	defer worker.Shutdown()

	callChan := make(chan any)
	worker.Register(queueName, func(c context.Context, received []byte) error {
		fmt.Println("recevied message", string(received))
		assert.Equal(t, payload, received)
		callChan <- true
		return nil
	})

	worker.Start()

	// sand random stuff make sure queueName actually works
	producer.Enqueue("otherqueue", []byte("to be ignored"))

	schedule := time.Now().UTC().Add(5 * time.Second)
	producer.Schedule(queueName, payload, schedule)

	assertSignal(t, time.Until(schedule), callChan, "handler not called")
}

func assertSignal(t *testing.T, timeout time.Duration, c <-chan any, msg string) {
	ti := time.NewTicker(timeout)
	defer ti.Stop()
	select {
	case <-c:
	case <-ti.C:
		require.FailNow(t, "no signal received: "+msg)
	}
}
