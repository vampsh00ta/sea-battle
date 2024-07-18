package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"seabattle/internal/entity"
)

//type BattleField interface {
//	GetBattleField(ctx context.Context, sessionID, idChatKey string, myField bool) (*entity.BattleField, error)
//	GetBattleFields(ctx context.Context, sessionID, idChatKey string) ([]*entity.BattleField, error)
//	SetBattleField(ctx context.Context, sessionID, idChatKey string, fields *entity.BattleField, myField bool) error
//}

func (db db) GetBattleField(ctx context.Context, sessionID, tgID string, myField bool) (*entity.BattleField, error) {
	// Агрегационный конвейер
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"users.tg_id", tgID}, {"session_id", sessionID}}}},
		{{"$project", bson.D{
			{"users", bson.D{
				{"$filter", bson.D{
					{"input", "$users"},
					{"as", "user"},
					{"cond", bson.D{{"$eq", bson.A{"$$user.tg_id", tgID}}}},
				}},
			}},
			{"_id", 0}, // Исключаем поле _id
		}}},
	}

	// Переменная для хранения результата
	var result struct {
		Users []entity.User `bson:"users"`
	}

	// Выполнение агрегации
	cursor, err := db.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		if len(result.Users) > 0 {
			if myField {
				return result.Users[0].MyField, nil
			} else {
				return result.Users[0].EnemyField, nil

			}
		}
	}

	return nil, errors.New("user not found")
}

func (db db) GetBattleFields(ctx context.Context, sessionID, idChatKey string) ([]*entity.BattleField, error) {
	var user entity.User

	filter := bson.D{{"session_id", sessionID}, {"users.tg_id", idChatKey}}
	if err := db.collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}
	if user.MyField != nil || user.EnemyField != nil {
		return nil, fmt.Errorf("user error")
	}
	return []*entity.BattleField{user.MyField, user.EnemyField}, nil
}
func (db db) SetBattleField(ctx context.Context, sessionID, idChatKey string, fields *entity.BattleField, myField bool) error {
	var fieldName string
	if myField {
		fieldName = "my_field"
	} else {
		fieldName = "enemy_field"
	}
	filter := bson.D{{"session_id", sessionID}, {"users.tg_id", idChatKey}}
	update := bson.M{
		"$set": bson.M{
			"users.$[elem]." + fieldName: fields,
		},
	}
	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem.tg_id": idChatKey}},
	})
	_, err := db.collection.UpdateOne(ctx, filter, update, arrayFilters)
	if err != nil {
		return err
	}

	return nil
}
