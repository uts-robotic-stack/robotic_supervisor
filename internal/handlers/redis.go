package handlers

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisHandler struct {
	client *redis.Client
}

func NewRedisHandler(addr, password string, db int) *RedisHandler {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisHandler{client: rdb}
}

func (h *RedisHandler) Set(ctx context.Context, key string, value interface{}) error {
	err := h.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (h *RedisHandler) Get(ctx context.Context, key string) (string, error) {
	val, err := h.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (h *RedisHandler) Delete(ctx context.Context, key string) error {
	err := h.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (h *RedisHandler) GetAllKeys(ctx context.Context) ([]string, error) {
	keys, err := h.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}
