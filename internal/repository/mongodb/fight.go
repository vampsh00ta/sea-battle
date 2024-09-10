package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	irep "seabattle/internal/app/repository/mongorep"
	"seabattle/internal/entity"
)

//	type Fight interface {
//		GetFight(ctx context.Context, sessionID string) (entity.Fight, error)
//		UpdateFight(ctx context.Context, fight entity.Fight) error
//		CreateFight(ctx context.Context, fight entity.Fight) error
//	}
type Fight struct {
	collection *mongo.Collection
}

func NewFight(collection *mongo.Collection) irep.Fight {
	return &Fight{collection: collection}
}
func (db Fight) CreateFight(ctx context.Context, fight entity.Fight) error {
	_, err := db.collection.InsertOne(ctx, fight)
	if err != nil {
		return err
	}
	return nil

}
func (db Fight) UpdateFight(ctx context.Context, fight entity.Fight) error {
	filter := bson.D{{"session_id", fight.SessionId}}
	_, err := db.collection.UpdateOne(ctx, filter, bson.D{{"$set", fight}})
	if err != nil {
		return err
	}

	return nil
}

func (db Fight) GetFight(ctx context.Context, sessionID string) (entity.Fight, error) {
	var fight entity.Fight
	filter := bson.D{{"session_id", sessionID}}
	if err := db.collection.FindOne(ctx, filter).Decode(&fight); err != nil {
		return entity.Fight{}, err
	}
	return fight, nil

}
