package redis

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"seabattle/internal/repository/models"
)

type Repository interface {
	GetBattleField(ctx context.Context, idChatKey string) (models.BattleField, error)
	SetBattleField(ctx context.Context, idChatKey string, fields string, myField bool) error

	GetSessionByChatId(ctx context.Context, idChatKey string) (string, error)
	CreateSessionByChatId(ctx context.Context, idChatKey1, idChatKey2 string) error
}

type Redis struct {
	client *redis.Client
}

func (r Redis) GetBattleField(ctx context.Context, idChatKey string) (models.BattleField, error) {
	//TODO implement me
	panic("implement me")
}

func (r Redis) CreateSessionByChatId(ctx context.Context, idChatKey1, idChatKey2 string) error {
	sessionId := uuid.New().String()
	var err error
	err = r.client.HSet(ctx, idChatKey1, models.User{SessionId: sessionId}).Err()
	if err != nil {
		return err
	}
	err = r.client.HSet(ctx, idChatKey2, models.User{SessionId: sessionId}).Err()

	if err != nil {
		return err
	}

	return nil
}

func (r Redis) GetSessionByChatId(ctx context.Context, idChatKey string) (string, error) {
	res := r.client.HGet(ctx, idChatKey, models.BattleSession)
	if res.Err() != nil {
		return "", nil
	}
	return res.Val(), nil
}
func (r Redis) SetBattleField(ctx context.Context, idChatKey string, fields string, myField bool) error {
	var err error
	var key string
	if myField {
		key = models.MyField
	} else {
		key = models.EnemyField
	}
	err = r.client.HSet(ctx, idChatKey, key, fields).Err()
	if err != nil {
		return err
	}
	return nil

}

func New(client *redis.Client) Repository {
	return &Redis{
		client,
	}
}
