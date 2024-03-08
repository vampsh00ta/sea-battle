package redis

import (
	"context"
	"seabattle/internal/repository/models"
)

type User interface {
	GetUserByChatId(ctx context.Context, idChatKey string) (models.User, error)
}

func (r Redis) GetUserByChatId(ctx context.Context, idChatKey string) (models.User, error) {
	var user models.User
	if err := r.client.HGetAll(ctx, idChatKey).Scan(&user); err != nil {
		return models.User{}, nil
	}
	return user, nil
}
