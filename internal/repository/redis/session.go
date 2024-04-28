package redis

import (
	"context"
	"github.com/google/uuid"
	"seabattle/internal/models"
)

type Session interface {
	GetSession(ctx context.Context, sessionId string) (models.Session, error)
	SetSession(ctx context.Context, sessionId string, session models.Session) error
	GetSessionByChatId(ctx context.Context, idChatKey string) (string, error)
	CreateSessionByChatId(ctx context.Context, idChatKey1, idChatKey2 string) (string, error)
	CreateSessionOnePerson(ctx context.Context, idChat string) (string, error)
}

func (r Redis) SetSession(ctx context.Context, sessionId string, session models.Session) error {
	if err := r.client.HSet(ctx, sessionId, session).Err(); err != nil {
		return err
	}
	return nil
}
func (r Redis) GetSession(ctx context.Context, sessionId string) (models.Session, error) {
	var session models.Session
	if err := r.client.HGetAll(ctx, sessionId).Scan(&session); err != nil {
		return models.Session{}, err
	}
	return session, nil
}

func (r Redis) CreateSessionOnePerson(ctx context.Context, idChat string) (string, error) {
	sessionId := battleSession + "_" + uuid.New().String()
	var err error
	session := models.Session{
		TgId1: idChat,
		TgId2: "",
		Ready: 1,
	}
	err = r.client.HSet(ctx, sessionId, session).Err()
	if err != nil {
		return "", err
	}

	return sessionId, nil
}
func (r Redis) CreateSessionByChatId(ctx context.Context, idChatKey1, idChatKey2 string) (string, error) {
	sessionId := battleSession + "_" + uuid.New().String()
	var err error
	session := models.Session{
		TgId1: idChatKey1,
		TgId2: idChatKey2,
		Ready: 0,
		Turn:  idChatKey1,
	}
	err = r.client.HSet(ctx, sessionId, session).Err()
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (r Redis) GetSessionByChatId(ctx context.Context, idChatKey string) (string, error) {
	res := r.client.HGet(ctx, idChatKey, battleSession)
	if res.Err() != nil {
		return "", nil
	}
	return res.Val(), nil
}
