package storage

import (
	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(address string, password string) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})
	return &Redis{client: client}
}

func (redis Redis) Set(key string, value interface{}) error {
	return redis.client.Set(key, value, 0).Err()
}

func (red Redis) Get(key string) (interface{}, error) {
	res, err := red.client.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return res, err
}

func (redis Redis) Delete(key string) error {
	return redis.client.Del(key).Err()
}
