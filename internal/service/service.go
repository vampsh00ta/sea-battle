package service

import (
	kafkago "github.com/segmentio/kafka-go"
	"seabattle/config"
	irep "seabattle/internal/app/repository"
	isrvc "seabattle/internal/app/service"
)

type service struct {
	mongo    irep.Repository
	gameConf *config.Game
	kafka    *kafkago.Writer
}

func New(mongo irep.Repository, gameConf *config.Game, kafka *kafkago.Writer) isrvc.Service {
	return &service{
		mongo:    mongo,
		gameConf: gameConf,
		kafka:    kafka}
}
