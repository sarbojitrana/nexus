package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	pgxzero "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/newrelic/go-agent/v3/integrations/nrpgx5"
	"github.com/rs/zerolog"
	"github.com/sarbojitrana/nexus/internal/config"
	loggerConfig "github.com/sarbojitrana/nexus/internal/logger"
)

type Database struct {
	Pool *pgxpool.Pool
	log  *zerolog.Logger
}

// multiTracer allows chaining multiple tracers
type multiTracer struct {
	tracers []any
}

const DatabasePingTimemout = 10

// TraceQueryStart implements pgx tracer interface

func (mt *multiTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, tracer := range mt.tracers {
		if t, ok := tracer.(interface {
			TraceQueryStart(context.Context, *pgx.Conn, pgx.TraceQueryStartData) context.Context
		}); ok {
			ctx = t.TraceQueryStart(ctx, conn, data)
		}
	}

	return ctx
}

// TraceQueryEnd implements pgx tracer interface
func (mt *multiTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, tracer := range mt.tracers {
		if t, ok := tracer.(interface {
			TraceQueryEnd(context.Context, *pgx.Conn, pgx.TraceQueryEndData)
		}); ok {
			t.TraceQueryEnd(ctx, conn, data)
		}
	}
}

// setup a new database connection

// loggerService is the application logger

func New(cfg *config.Config, logger *zerolog.Logger, loggerService *loggerConfig.LoggerService) (*Database, error) {
	hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port)) // to make hostPort as host:Port, should not be done explicitely by string manipulation

	//URL-encode the password

	encodedPassword := url.QueryEscape(cfg.Database.Password)

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		encodedPassword,
		hostPort,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse pgx pool config: %w", err)
	}

	if loggerService != nil && loggerService.GetApplication() != nil {
		pgxPoolConfig.ConnConfig.Tracer = nrpgx5.NewTracer() //whenever a database event happens call this tracer on new relic if its running
	}

	if cfg.Primary.Env == "local" {
		globalLevel := logger.GetLevel()
		pgxLogger := loggerConfig.NewPgxLogger(globalLevel)

		// Chain tracers : New relic first, then local logging
		// tracelog is aimed for SQL queries which it logs whereas new relic tracer  creates APM spans and sends data to new relic which in turn will help in observability
		if pgxPoolConfig.ConnConfig.Tracer != nil { // if new relic tracer exists, create a multitracer one for new relic and other for local, new relic tracer will have higher priority

			localTracer := &tracelog.TraceLog{ // when pgx executes a query Tracelog reacts by logging SQL statement, argument, duration and errors and TraceQueryStart and end are called
				Logger:   pgxzero.NewLogger(pgxLogger),
				LogLevel: tracelog.LogLevel(loggerConfig.GetPgxTraceLogLevel(globalLevel)),
			}

			pgxPoolConfig.ConnConfig.Tracer = &multiTracer{
				tracers: []any{pgxPoolConfig.ConnConfig.Tracer, localTracer},
			}
		} else {
			pgxPoolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
				Logger:   pgxzero.NewLogger(pgxLogger),
				LogLevel: tracelog.LogLevel(loggerConfig.GetPgxTraceLogLevel(globalLevel)),
			}
		}
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig) // make a pool with new context

	if err != nil {
		return nil, fmt.Errorf("failed to creaet pgx pool : %w", err)
	}

	database := &Database{
		Pool: pool,
		log:  logger,
	}

	ctx, cancel := context.WithTimeout(context.Background(), DatabasePingTimemout*time.Second)
	defer cancel()
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping to database : %w", err)
	}

	logger.Info().Msg("Connected to the database")
	return database, nil
}

func (db *Database) Close() error {
	db.log.Info().Msg("Closing database connection pool")
	db.Pool.Close()
	return nil
}
