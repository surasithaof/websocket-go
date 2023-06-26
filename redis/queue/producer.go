package redisqueues

import (
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
)

type Producer struct {
	client *asynq.Client
}

func NewProducer(config ClientConfig) QueueProducer {
	opt, err := redis.ParseURL(config.URL)
	if err != nil {
		panic(err)
	}
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: opt.Addr,
		// Username: opt.Username,
		Password: opt.Password,
		DB:       opt.DB,
		// DialTimeout:  config.DialTimeout,
		// ReadTimeout:  config.ReadTimeout,
		// WriteTimeout: config.WriteTimeout,
		// PoolSize:     config.PoolSize,
	})
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
