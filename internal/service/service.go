package service

import (
	kafkago "github.com/segmentio/kafka-go"
	"seabattle/config"
	"seabattle/internal/repository/psql"
	"seabattle/internal/repository/redis"
	action "seabattle/internal/service/action"
)

type Service interface {
	Fight
	Session
	CodeGenerator
}

type service struct {
	redis    redis.Repository
	psql     psql.Repository
	gameConf *config.Game
	action   action.Action
	kafka    *kafkago.Writer
}

func New(repo redis.Repository, psql psql.Repository, gameConf *config.Game, kafka *kafkago.Writer) Service {
	act := action.New(gameConf)
	return &service{redis: repo, psql: psql, gameConf: gameConf, action: act, kafka: kafka}
}
