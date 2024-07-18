package repository

import (
	"context"
	"seabattle/internal/entity"
)

type Fight interface {
	GetFight(ctx context.Context, sessionID string) (entity.Fight, error)
	UpdateFight(ctx context.Context, fight entity.Fight) error
	CreateFight(ctx context.Context, fight entity.Fight) error
}
