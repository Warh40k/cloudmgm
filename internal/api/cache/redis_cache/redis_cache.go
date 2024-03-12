package redis_cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

const TTL = 12 * time.Hour

func NewRedisConn(ctx context.Context, cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	status := client.Ping(ctx)
	if err := status.Err(); err != nil {
		return nil, err
	}
	return client, nil
}
