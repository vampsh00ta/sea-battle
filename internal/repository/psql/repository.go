package psql

import (
	"context"
	"seabattle/pkg/client"
)

type Db struct {
	client client.Client
}

type Repository interface {
	GetSessionByCode(ctx context.Context, code string) (string, error)
	AddSession(ctx context.Context, code, sessionId string) error
}

func (db Db) GetSessionByCode(ctx context.Context, code string) (string, error) {
	var sessionId string
	q := `select session from link where code = $1`

	if err := db.client.QueryRow(ctx, q, code).Scan(&sessionId); err != nil {
		return "", err
	}

	return sessionId, nil
}
func (db Db) AddSession(ctx context.Context, code, sessionId string) error {
	q := `insert into link(code,session) values($1,$2) returning code`

	if err := db.client.QueryRow(ctx, q, code, sessionId).Scan(&code); err != nil {
		return err
	}

	return nil
}
func New(client client.Client) Repository {

	return &Db{
		client: client,
	}
}
