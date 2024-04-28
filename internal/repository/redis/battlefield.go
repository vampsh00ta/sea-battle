package redis

import (
	"context"
	"encoding/json"
	"seabattle/internal/models"
	"strconv"
)

type BattleField interface {
	GetBattleField(ctx context.Context, idChatKey string, myField bool) (*models.BattleField, error)
	GetBattleFields(ctx context.Context, idChatKey string) ([]*models.BattleField, error)
	SetBattleField(ctx context.Context, idChatKey string, fields *models.BattleField, myField bool) error
	SetPoint(ctx context.Context, idChatKey string, x, y int) error
	GetPoint(ctx context.Context, idChatKey string) (int, int, error)
}

func (r Redis) GetPoint(ctx context.Context, idChatKey string) (int, int, error) {
	x_redis := r.client.HGet(ctx, idChatKey, "curr_x")
	if x_redis.Err() != nil {
		return -1, -1, x_redis.Err()
	}
	xStr := x_redis.Val()
	x, _ := strconv.Atoi(xStr)
	y_redis := r.client.HGet(ctx, idChatKey, "curr_y")
	if x_redis.Err() != nil {
		return -1, -1, x_redis.Err()
	}
	yStr := y_redis.Val()
	y, _ := strconv.Atoi(yStr)
	return x, y, nil
}
func (r Redis) SetPoint(ctx context.Context, idChatKey string, x, y int) error {
	err := r.client.HSet(ctx, idChatKey, "curr_x", x, "curr_y", y).Err()
	if err != nil {
		return err
	}
	return nil
}
func (r Redis) GetBattleFields(ctx context.Context, idChatKey string) ([]*models.BattleField, error) {
	my, err := r.GetBattleField(ctx, idChatKey, true)
	if err != nil {
		return nil, err
	}
	enemy, err := r.GetBattleField(ctx, idChatKey, false)
	if err != nil {
		return nil, err
	}
	return []*models.BattleField{my, enemy}, nil
}

func (r Redis) GetBattleField(ctx context.Context, idChatKey string, me bool) (*models.BattleField, error) {
	var key string
	if me {
		key = myField
	} else {
		key = enemyField
	}

	res := r.client.HGet(ctx, idChatKey, key)
	if res.Err() != nil {
		return nil, res.Err()
	}
	dataStr := res.Val()
	var data models.BattleField
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return nil, err
	}
	return &data, err
}
func (r Redis) SetBattleField(ctx context.Context, idChatKey string, fields *models.BattleField, me bool) error {
	var key string
	if me {
		key = myField
	} else {
		key = enemyField
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
