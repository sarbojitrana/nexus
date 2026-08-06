package job

import (
	"strings"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/sarbojitrana/nexus/internal/config"
)

// parseRedisConnOpt accepts either a bare "host:port" (local dev,
// docker-compose) or a full "redis://[:password@]host:port" URL (Render's
// Key Value connectionString) -- asynq.RedisClientOpt.Addr only understands
// the former, so a full URL needs to go through ParseRedisURI instead.
func parseRedisConnOpt(address string, logger *zerolog.Logger) asynq.RedisConnOpt {
	if strings.Contains(address, "://") {
		connOpt, err := asynq.ParseRedisURI(address)
		if err != nil {
			logger.Fatal().Err(err).Msg("invalid redis address")
		}
		return connOpt
	}
	return asynq.RedisClientOpt{Addr: address}
}

type JobService struct { // background task queue
	Client *asynq.Client
	server *asynq.Server
	logger *zerolog.Logger
}

func NewJobService(logger *zerolog.Logger, cfg *config.Config) *JobService {
	redisConnOpt := parseRedisConnOpt(cfg.Redis.Address, logger)

	client := asynq.NewClient(redisConnOpt) // will only enqueue tasks

	server := asynq.NewServer( // runs on goroutine
		redisConnOpt,
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
	mux.HandleFunc(TaskSignIn, j.handleSignInEmailTask)
	mux.HandleFunc(TaskPasswordChanged, j.handlePasswordChangedEmailTask)
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
