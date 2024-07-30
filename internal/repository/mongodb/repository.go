package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	irep "seabattle/internal/app/repository/mongorep"
)

type db struct {
	collection *mongo.Collection
}

func New(collection *mongo.Collection) irep.Repository {
	return db{collection}
}
