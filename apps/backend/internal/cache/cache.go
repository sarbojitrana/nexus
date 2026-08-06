package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Client struct {
	redis  *redis.Client
	logger *zerolog.Logger
}

func NewClient(redis *redis.Client, logger *zerolog.Logger) *Client {
	return &Client{redis: redis, logger: logger}
}

func (c *Client) Get(ctx context.Context, key string, dest any) bool {
	raw, err := c.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err != redis.Nil {
			c.logger.Warn().Err(err).Str("key", key).Msg("cache get failed, falling back to source")
		}
		return false
	}

	if err := json.Unmarshal(raw, dest); err != nil {
		c.logger.Warn().Err(err).Str("key", key).Msg("cache value corrupt, falling back to source")
		return false
	}

	return true
}

func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) {
	raw, err := json.Marshal(value)
	if err != nil {
		c.logger.Warn().Err(err).Str("key", key).Msg("failed to marshal value for cache")
		return
	}

	if err := c.redis.Set(ctx, key, raw, ttl).Err(); err != nil {
		c.logger.Warn().Err(err).Str("key", key).Msg("cache set failed")
	}
}

func (c *Client) Delete(ctx context.Context, keys ...string) {
	if len(keys) == 0 {
		return
	}
	if err := c.redis.Del(ctx, keys...).Err(); err != nil {
		c.logger.Warn().Err(err).Strs("keys", keys).Msg("cache delete failed")
	}
}
