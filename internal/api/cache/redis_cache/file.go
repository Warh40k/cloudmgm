package redis_cache

import (
	"context"
	"encoding/json"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type FileCache struct {
	ctx   context.Context
	cache *redis.Client
}

func NewFileCache(ctx context.Context, cache *redis.Client) *FileCache {
	return &FileCache{ctx: ctx, cache: cache}
}

func (c *FileCache) GetFileInfo(id uuid.UUID) (domain.File, error) {
	var file domain.File
	strfile, err := c.cache.Get(c.ctx, id.String()).Bytes()
	if err != nil {
		return file, err
	}
	err = json.Unmarshal(strfile, &file)
	if err != nil {
		return file, err
	}

	return file, nil
}

func (c *FileCache) SetFileInfo(file domain.File) error {
	jfile, err := json.Marshal(file)
	if err != nil {
		return err
	}
	err = c.cache.Set(c.ctx, file.Id.String(), jfile, TTL).Err()
	if err != nil {
		return err
	}
	return nil
}
