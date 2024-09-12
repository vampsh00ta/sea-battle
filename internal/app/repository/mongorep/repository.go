package mongorep

import (
	"context"
	"seabattle/internal/entity"
)

//	type Repository interface {
//		User
//		Fight
//		Field
//	}
type BattleField interface {
	GetBySessionID(ctx context.Context, sessionID, idChatKey string, myField bool) (*entity.BattleField, error)
	GetAll(ctx context.Context, sessionID, idChatKey string) ([]*entity.BattleField, error)
	Set(ctx context.Context, sessionID, idChatKey string, fields *entity.BattleField, myField bool) error
}
type Fight interface {
	GetBySessionID(ctx context.Context, sessionID string) (entity.Fight, error)
	Update(ctx context.Context, fight entity.Fight) error
	Create(ctx context.Context, fight entity.Fight) error
}
type User interface {
	Set(ctx context.Context, sessionID string, user entity.User) error
	GetByTgID(ctx context.Context, tgID string) (entity.User, error)
	SetFieldQueryId(ctx context.Context, sessionID, tgId, queryID string, my bool) error
	SetPoint(ctx context.Context, sessionID, idChatKey string, point entity.Point) error
	GetPoint(ctx context.Context, sessionID, idChatKey string) (entity.Point, error)
}
