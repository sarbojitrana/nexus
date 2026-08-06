package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/sarbojitrana/nexus/internal/cache"
	"github.com/sarbojitrana/nexus/internal/config"
	"github.com/sarbojitrana/nexus/internal/database"
	"github.com/sarbojitrana/nexus/internal/lib/job"
	"github.com/sarbojitrana/nexus/internal/lib/ws"
	loggerPkg "github.com/sarbojitrana/nexus/internal/logger"
	"github.com/sarbojitrana/nexus/internal/search"
	"github.com/sarbojitrana/nexus/internal/storage"
)

// parseRedisAddress accepts either a bare "host:port" (local dev,
// docker-compose) or a full "redis://[:password@]host:port" URL (Render's
// Key Value connectionString) -- go-redis's Options.Addr only understands
// the former, so a full URL needs to go through ParseURL instead.
func parseRedisAddress(address string) (*redis.Options, error) {
	if strings.Contains(address, "://") {
		return redis.ParseURL(address)
	}
	return &redis.Options{Addr: address}, nil
}

type Server struct {
	Config        *config.Config
	Logger        *zerolog.Logger
	LoggerService *loggerPkg.LoggerService
	DB            *database.Database
	Redis         *redis.Client
	httpServer    *http.Server
	Job           *job.JobService
	Search        *search.Client
	Storage       *storage.Client
	Hub           *ws.Hub
	Cache         *cache.Client
}

func New(cfg *config.Config, logger *zerolog.Logger, loggerService *loggerPkg.LoggerService) (*Server, error) {
	db, err := database.New(cfg, logger, loggerService)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	redisOpts, err := parseRedisAddress(cfg.Redis.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis address: %w", err)
	}
	redisClient := redis.NewClient(redisOpts)

	if loggerService != nil && loggerService.GetApplication() != nil {
		redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error().Err(err).Msg("failed to connect to Redis, continuing without Redis")
	}

	jobService := job.NewJobService(logger, cfg)
	jobService.InitHandlers(cfg, logger)

	if err := jobService.Start(); err != nil {
		return nil, err
	}

	initCtx, initCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer initCancel()

	searchClient := search.NewClient(initCtx, cfg, logger)

	storageClient, err := storage.NewClient(initCtx, cfg)
	if err != nil {
		logger.Error().Err(err).Msg("failed to initialize storage client, uploads disabled")
		storageClient = nil
	}

	server := &Server{
		Config:        cfg,
		Logger:        logger,
		LoggerService: loggerService,
		DB:            db,
		Redis:         redisClient,
		Job:           jobService,
		Search:        searchClient,
		Storage:       storageClient,
		Hub:           ws.NewHub(logger),
		Cache:         cache.NewClient(redisClient, logger),
	}

	return server, nil

}

func (s *Server) SetupHTTPServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:         ":" + s.Config.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeout) * time.Second,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}

	s.Logger.Info().
		Str("port", s.Config.Server.Port).
		Str("env", s.Config.Primary.Env).
		Msg("starting server")

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server : %w", err)
	}

	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("Failed to close database connection : %w", err)
	}

	if s.Job != nil {
		s.Job.Stop()
	}

	return nil
}
