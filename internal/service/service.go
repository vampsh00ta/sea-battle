package service

import (
	"seabattle/internal/repository/psql"
	"seabattle/internal/repository/redis"
)

type Service interface {
	Fight
	Session
	CodeGenerator
}

type service struct {
	redis redis.Repository
	psql  psql.Repository
}

func New(repo redis.Repository, psql psql.Repository) Service {
	return &service{redis: repo, psql: psql}
}
