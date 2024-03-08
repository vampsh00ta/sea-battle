package redis

import (
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	User
	Fight
	Session
	BattleField
}

type Redis struct {
	client *redis.Client
}

const (
	battleSession = "tg_battle_session"
	myField       = "my_field"
	enemyField    = "enemy_field"
)

func New(client *redis.Client) Repository {
	return &Redis{
		client,
	}
}
