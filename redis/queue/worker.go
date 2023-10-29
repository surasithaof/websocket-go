package redisqueues

import (
	"context"
	"strings"

	"github.com/hibiken/asynq"
)

func NewWorker(conf WorkerConfig) QueueWorker {
	var logLevel asynq.LogLevel
	logLevel.Set(conf.LogLevel)

	queueConf := asynq.Config{
		Concurrency:              conf.Concurrency,
		StrictPriority:           conf.StrictPriority,
		ShutdownTimeout:          conf.ShutdownTimeout,
		HealthCheckInterval:      conf.HealthCheckInterval,
		DelayedTaskCheckInterval: conf.DelayedTaskCheckInterval,
		GroupGracePeriod:         conf.GroupGracePeriod,
		GroupMaxDelay:            conf.GroupMaxDelay,
		GroupMaxSize:             conf.GroupMaxSize,
		// Logger:                   monitoring.Logger().Sugar(),
		LogLevel: logLevel,
	}

	var clientOpt asynq.RedisConnOpt
	if conf.ClusterModeEnable {
		clientOpt = asynq.RedisClusterClientOpt{
			Addrs: strings.Split(conf.URL, ","),
		}
	} else {
		clientOpt = asynq.RedisClientOpt{
			Addr: conf.URL,
		}
	}

	srv := asynq.NewServer(
		clientOpt,
		queueConf,
	)

	w := &Worker{
		handlers:    map[string]func(context.Context, *asynq.Task) error{},
		asyncServer: srv,
	}

	return w
}

type Worker struct {
	handlers    map[string]func(context.Context, *asynq.Task) error
	asyncServer *asynq.Server
}

func (w *Worker) Register(queueName string, handler func(c context.Context, payload []byte) error) {
	w.handlers[queueName] = func(c context.Context, t *asynq.Task) error {
		err := handler(c, t.Payload())
		return err
	}
}

func (w *Worker) Start() {
	// mux maps a type to a handler
	mux := asynq.NewServeMux()

	for pattern, handler := range w.handlers {
		mux.HandleFunc(pattern, handler)
	}
	w.asyncServer.Start(mux)

	// if err := srv.Run(mux); err != nil {
	// 	log.Fatalf("could not run server: %v", err)
	// }
}

func (w *Worker) Shutdown() {
	w.asyncServer.Shutdown()
}
