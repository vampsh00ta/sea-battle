package repository

import (
	"context"
	"seabattle/internal/entity"
)

type Field interface {
	GetBattleField(ctx context.Context, sessionID, idChatKey string, myField bool) (*entity.BattleField, error)
	GetBattleFields(ctx context.Context, sessionID, idChatKey string) ([]*entity.BattleField, error)
	SetBattleField(ctx context.Context, sessionID, idChatKey string, fields *entity.BattleField, myField bool) error
}
