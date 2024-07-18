package repository

import (
	"context"
	"seabattle/internal/entity"
)

type User interface {
	GetUser(ctx context.Context, tgID string) (entity.User, error)
	SetFieldQueryId(ctx context.Context, sessionID, tgId, queryID string, my bool) error
	SetPoint(ctx context.Context, sessionID, idChatKey string, point entity.Point) error
	GetPoint(ctx context.Context, sessionID, idChatKey string) (entity.Point, error)
}
