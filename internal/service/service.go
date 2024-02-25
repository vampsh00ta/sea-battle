package service

import "seabattle/internal/repository/redis"

type Service interface {
}

type service struct {
	repo redis.Repository
}

func New(repo redis.Repository) Service {
	return &service{
		repo: repo,
	}
}
