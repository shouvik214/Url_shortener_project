package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func (r *RedisRepository) Close() error {
	return r.client.Close()
}

func NewRedisRepository(
	addr string,
) *RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) Set(
	ctx context.Context,
	key string,
	value string,
	expiration time.Duration,
) error {
	return r.client.Set(
		ctx,
		key,
		value,
		expiration,
	).Err()
}

func (r *RedisRepository) Get(
	ctx context.Context,
	key string,
) (string, error) {
	value, err := r.client.Get(ctx, key).Result()

	if err == redis.Nil {
		// Cache miss — this is not a Redis failure.
		return "", nil
	}

	if err != nil {
		// Redis itself is unavailable.
		return "", err
	}

	return value, nil
}

func (r *RedisRepository) Delete(
	ctx context.Context,
	key string,
) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisRepository) Ping(
	ctx context.Context,
) error {
	return r.client.Ping(ctx).Err()
}
