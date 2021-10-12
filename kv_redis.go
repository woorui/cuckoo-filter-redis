package filter

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

type RedisKV struct {
	key string
	rdb *redis.Client
}

func NewRedisKV(rdb *redis.Client, key string) *RedisKV {
	return &RedisKV{rdb: rdb, key: key}
}

func (kv *RedisKV) Get(ctx context.Context) ([]byte, error) {
	b, se := kv.rdb.Get(ctx, kv.key).Bytes()
	if se != nil {
		if errors.Is(se, redis.Nil) {
			return []byte{}, nil
		}
		return []byte{}, se
	}
	return b, nil
}

func (kv *RedisKV) Set(ctx context.Context, value []byte) error {
	return kv.rdb.Set(ctx, kv.key, value, 0).Err()
}
