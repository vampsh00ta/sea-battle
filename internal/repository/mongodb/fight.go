package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"seabattle/internal/entity"
)

//type Fight interface {
//	GetFight(ctx context.Context, sessionID string) (entity.Fight, error)
//	UpdateFight(ctx context.Context, fight entity.Fight) error
//	CreateFight(ctx context.Context, fight entity.Fight) error
//}

func (db db) CreateFight(ctx context.Context, fight entity.Fight) error {
	_, err := db.collection.InsertOne(ctx, fight)
	if err != nil {
		return err
	}
	return nil

}
func (db db) UpdateFight(ctx context.Context, fight entity.Fight) error {
	filter := bson.D{{"session_id", fight.SessionId}}
	_, err := db.collection.UpdateOne(ctx, filter, bson.D{{"$set", fight}})
	if err != nil {
		return err
	}

	return nil
}

func (db db) GetFight(ctx context.Context, sessionID string) (entity.Fight, error) {
	var fight entity.Fight
	filter := bson.D{{"session_id", sessionID}}
	if err := db.collection.FindOne(ctx, filter).Decode(&fight); err != nil {
		return entity.Fight{}, err
	}
	return fight, nil

}
