package redis

import (
	"context"
	"encoding/json"
	"seabattle/internal/repository/models"
)

type BattleField interface {
	GetBattleField(ctx context.Context, idChatKey string, myField bool) (*models.BattleField, error)
	GetBattleFields(ctx context.Context, idChatKey string) ([]*models.BattleField, error)
	SetBattleField(ctx context.Context, idChatKey string, fields *models.BattleField, myField bool) error
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
	var err error
	var key string
	if me {
		key = myField
	} else {
		key = enemyField
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
func (r Redis) SetBattleField(ctx context.Context, idChatKey string, fields *models.BattleField, me bool) error {
	var err error
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
