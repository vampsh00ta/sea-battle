package redis

import "github.com/redis/go-redis/v9"

type RedisRep interface {
}
type Redis struct {
	client *redis.Client
}

func New(client *redis.Client) RedisRep {
	return &Redis{
		client,
	}
}
