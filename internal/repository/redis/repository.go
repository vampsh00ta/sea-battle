package redis

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"seabattle/internal/repository/models"
)

type Repository interface {
	GetBattleField(ctx context.Context, idChatKey string, myField bool) (*models.BattleField, error)
	SetBattleField(ctx context.Context, idChatKey string, fields *models.BattleField, myField bool) error

	GetUserByChatId(ctx context.Context, idChatKey string) (models.User, error)

	GetSessionByChatId(ctx context.Context, idChatKey string) (string, error)
	CreateSessionByChatId(ctx context.Context, idChatKey1, idChatKey2 string) (string, error)
}

type Redis struct {
	client *redis.Client
}

func (r Redis) GetUserByChatId(ctx context.Context, idChatKey string) (models.User, error) {
	var user models.User
	if err := r.client.HGetAll(ctx, idChatKey).Scan(&user); err != nil {
		return models.User{}, nil
	}
	return user, nil
}
func (r Redis) CreateSessionByChatId(ctx context.Context, idChatKey1, idChatKey2 string) (string, error) {
	sessionId := models.BattleSession + "_" + uuid.New().String()
	var err error
	err = r.client.HSet(ctx, sessionId,
		models.Session{
			TgId1: idChatKey1,
			TgId2: idChatKey2,
			Ready: 0,
			Stage: models.StagePicking,
		}).Err()
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (r Redis) GetSessionByChatId(ctx context.Context, idChatKey string) (string, error) {
	res := r.client.HGet(ctx, idChatKey, models.BattleSession)
	if res.Err() != nil {
		return "", nil
	}
	return res.Val(), nil
}

func (r Redis) GetBattleField(ctx context.Context, idChatKey string, myField bool) (*models.BattleField, error) {
	var err error
	var key string
	if myField {
		key = models.MyField
	} else {
		key = models.EnemyField
	}
	if err != nil {
		return nil, err
	}

	res := r.client.HGet(ctx, idChatKey, key)
	if res.Err() != nil {
		return nil, err
	}
	dataStr := res.Val()
	var data models.BattleField
	err = json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return nil, err
	}
	return &data, err
}
func (r Redis) SetBattleField(ctx context.Context, idChatKey string, fields *models.BattleField, myField bool) error {
	var err error
	var key string
	if myField {
		key = models.MyField
	} else {
		key = models.EnemyField
	}
	data, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	err = r.client.HSet(ctx, idChatKey, key, string(data)).Err()
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
