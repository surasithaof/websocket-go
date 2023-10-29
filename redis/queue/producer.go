package redisqueues

import (
	"log"
	"strings"
	"time"

	"github.com/hibiken/asynq"
)

type Producer struct {
	client *asynq.Client
}

func NewProducer(config Config) QueueProducer {
	var clientOpt asynq.RedisConnOpt
	if config.ClusterModeEnable {
		clientOpt = asynq.RedisClusterClientOpt{
			Addrs: strings.Split(config.URL, ","),
		}
	} else {
		clientOpt = asynq.RedisClientOpt{
			Addr: config.URL,
		}
	}

	client := asynq.NewClient(clientOpt)

	return NewProducerWithClient(client)
}

func NewProducerWithClient(client *asynq.Client) QueueProducer {
	producer := &Producer{client}
	return producer
}

func (p *Producer) Enqueue(jobType string, payload []byte) {

	task := asynq.NewTask(jobType, payload)

	info, err := p.client.Enqueue(task)
	if err != nil {
		log.Fatalf("could enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}

func (p *Producer) Schedule(jobType string, payload []byte, t time.Time) {

	task := asynq.NewTask(jobType, payload)

	info, err := p.client.Enqueue(task, asynq.ProcessAt(t))
	if err != nil {
		log.Fatalf("could not schedule task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}

func (p *Producer) Shutdown() error {
	return p.client.Close()
}
