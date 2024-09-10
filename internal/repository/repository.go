package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	irep "seabattle/internal/app/repository/mongorep"
	"seabattle/internal/repository/mongodb"
)

//type db struct {
//	collection *mongo.Collection
//}

type db struct {
	irep.User
	irep.Fight
	irep.Field
}

func New(collection *mongo.Collection) irep.Repository {
	return &db{
		User:  mongodb.NewUser(collection),
		Fight: mongodb.NewFight(collection),
		Field: mongodb.NewField(collection),
	}
}
