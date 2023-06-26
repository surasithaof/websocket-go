package redisqueues

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
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
		// Logger:                   log.New("worker"),
		LogLevel: logLevel,
	}

	// TODO cluster option
	// asynq api has pretty poor DX

	opt, err := redis.ParseURL(conf.URL)
	if err != nil {
		panic(err)
	}

	redisOpts := asynq.RedisClientOpt{
		Addr: opt.Addr,
		// Username:     opt.Username,
		Password:     opt.Password,
		DB:           opt.DB,
		DialTimeout:  opt.DialTimeout,
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,
		PoolSize:     opt.PoolSize,
	}

	srv := asynq.NewServer(
		redisOpts,
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

type WorkerConfig struct {
	// ClientConfig
	URL string `envconfig:"REDIS_URL" required:"true"`

	// Maximum number of concurrent processing of tasks.
	//
	// If set to a zero or negative value, NewServer will overwrite the value
	// to the number of CPUs usable by the current process.
	Concurrency int

	// StrictPriority indicates whether the queue priority should be treated strictly.
	//
	// If set to true, tasks in the queue with the highest priority is processed first.
	// The tasks in lower priority queues are processed only when those queues with
	// higher priorities are empty.
	StrictPriority bool

	// LogLevel specifies the minimum log level to enable.
	//
	// If unset, InfoLevel is used by default.
	LogLevel string `envconfig:"REDIS_QUEUES_LOG_LEVEL" default:"info"`

	// ShutdownTimeout specifies the duration to wait to let workers finish their tasks
	// before forcing them to abort when stopping the server.
	//
	// If unset or zero, default timeout of 8 seconds is used.
	ShutdownTimeout time.Duration `envconfig:"REDIS_QUEUES_SHUTDOWN_TIMEOUT" default:"8s"`

	// HealthCheckInterval specifies the interval between healthchecks.
	HealthCheckInterval time.Duration `envconfig:"REDIS_QUEUES_SHUTDOWN_TIMEOUT" default:"15s"`

	// DelayedTaskCheckInterval specifies the interval between checks run on 'scheduled' and 'retry'
	// tasks, and forwarding them to 'pending' state if they are ready to be processed.
	DelayedTaskCheckInterval time.Duration `envconfig:"REDIS_QUEUES_DELAYED_TASK_CHECK_INTERVAL" default:"5s"`

	// GroupGracePeriod specifies the amount of time the server will wait for an incoming task before aggregating
	// the tasks in a group. If an incoming task is received within this period, the server will wait for another
	// period of the same length, up to GroupMaxDelay if specified.
	//
	// If unset or zero, the grace period is set to 1 minute.
	// Minimum duration for GroupGracePeriod is 1 second. If value specified is less than a second, the call to
	// NewServer will panic.
	GroupGracePeriod time.Duration `envconfig:"REDIS_QUEUES_GROUP_GRACE_PERIOD" default:"1m"`

	// GroupMaxDelay specifies the maximum amount of time the server will wait for incoming tasks before aggregating
	// the tasks in a group.
	GroupMaxDelay time.Duration `envconfig:"REDIS_QUEUES_GROUP_MAX_DELAY" default:"0"`

	// GroupMaxSize specifies the maximum number of tasks that can be aggregated into a single task within a group.
	// If GroupMaxSize is reached, the server will aggregate the tasks into one immediately.
	//
	// If unset or zero, no size limit is used.
	GroupMaxSize int `envconfig:"REDIS_QUEUES_GROUP_MAX_SIZE" default:"0"`
}
