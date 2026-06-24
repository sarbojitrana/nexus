package job

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/sarbojitrana/nexus/internal/config"
)

type JobService struct { // background task queue
	Client *asynq.Client
	server *asynq.Server
	logger *zerolog.Logger
}

// using dependency injection

func NewJobService(logger *zerolog.Logger, cfg *config.Config) *JobService {
	redisAddr := cfg.Redis.Address

	client := asynq.NewClient(asynq.RedisClientOpt{ // will only enqueue tasks
		Addr: redisAddr,
	})

	server := asynq.NewServer( // runs on goroutine
		asynq.RedisClientOpt{
			Addr: redisAddr,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	return &JobService{
		Client: client,
		server: server,
		logger: logger,
	}
}

func (j *JobService) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, j.handleWelcomeEmailTask) // when a task TaskWelcome is pulled from Redis, call handleWelcomeEmailTask
	j.logger.Info().Msg("Starting background job server")
	if err := j.server.Start(mux); err != nil {
		return err
	}
	return nil
}

func (j *JobService) Stop() {
	j.logger.Info().Msg("Stopping background job server")
	j.server.Shutdown()
	j.Client.Close()
}
