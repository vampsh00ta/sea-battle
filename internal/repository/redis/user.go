package redis

import (
	"context"
	"encoding/json"
	"seabattle/internal/repository/models"
)

type User interface {
	GetUserByChatId(ctx context.Context, idChatKey string) (models.User, error)
	SetUser(ctx context.Context, user models.User) error
	SetFieldQueryId(ctx context.Context, tgId, queryId string, my bool) error
}

func (r Redis) SetUser(ctx context.Context, user models.User) error {

	my, err := json.Marshal(user.MyField)
	if err != nil {
		return err
	}
	enemy, err := json.Marshal(user.EnemyField)
	if err != nil {
		return err
	}
	userRedis := models.UserRedis{
		MyField:           string(my),
		EnemyField:        string(enemy),
		CurrX:             user.CurrX,
		CurrY:             user.CurrY,
		MyFieldQueryId:    user.MyFieldQueryId,
		EnemyFieldQueryId: user.EnemyFieldQueryId,
	}
	if err := r.client.HSet(ctx, user.TgId, userRedis); err != nil {
		return nil
	}
	return nil
}
func (r Redis) GetUser(ctx context.Context, tgId string) (models.User, error) {
	var user models.UserRedis
	if err := r.client.HGetAll(ctx, tgId).Scan(&user); err != nil {
		return models.User{}, err
	}
	var my, enemy models.BattleField
	if err := json.Unmarshal([]byte(user.MyField), &my); err != nil {
		return models.User{}, err
	}
	if err := json.Unmarshal([]byte(user.EnemyField), &enemy); err != nil {
		return models.User{}, err
	}
	userModel := models.User{
		TgId:              tgId,
		MyField:           &my,
		EnemyField:        &enemy,
		CurrX:             user.CurrX,
		CurrY:             user.CurrY,
		MyFieldQueryId:    user.MyFieldQueryId,
		EnemyFieldQueryId: user.EnemyFieldQueryId,
	}
	return userModel, nil
}
func (r Redis) GetUserByChatId(ctx context.Context, idChatKey string) (models.User, error) {
	var user models.User
	if err := r.client.HGetAll(ctx, idChatKey).Scan(&user); err != nil {
		return models.User{}, nil
	}
	return user, nil
}

func (r Redis) SetFieldQueryId(ctx context.Context, tgId, queryId string, my bool) error {

	var queryField string

	switch my {
	case true:
		queryField = models.MyFieldQueryId
	case false:
		queryField = models.EnemyFieldQueryId
	}
	if err := r.client.HSet(ctx, tgId, queryField, queryId).Err(); err != nil {
		return err
	}
	return nil
}
