package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	Client *redis.Client
}

func NewRedisStorage(addr, password string, db int) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisStorage{
		Client: client,
	}, nil
}

func (r *RedisStorage) Increment(ctx context.Context, key string, expiration int64) (int64, error) {
	countKey := "count:" + key

	pipe := r.Client.TxPipeline()
	incr := pipe.Incr(ctx, countKey)

	if expiration > 0 {
		pipe.Expire(ctx, countKey, time.Duration(expiration)*time.Second)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return incr.Val(), nil
}

func (r *RedisStorage) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisStorage) Set(ctx context.Context, key string, value interface{}, expiration int64) error {
	return r.Client.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
}

func (r *RedisStorage) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *RedisStorage) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	return result == 1, err
}

func (r *RedisStorage) Close() error {
	return r.Client.Close()
}
