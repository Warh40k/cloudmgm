package cache

import (
	"context"
	"github.com/Warh40k/cloud-manager/internal/api/cache/redis_cache"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type File interface {
	GetFileInfo(id uuid.UUID) (domain.File, error)
	SetFileInfo(file domain.File) error
}

type Cache struct {
	File
	ctx context.Context
}

func NewCache(ctx context.Context, client *redis.Client) *Cache {
	return &Cache{
		File: redis_cache.NewFileCache(ctx, client),
	}
}
