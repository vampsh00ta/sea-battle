package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	irep "seabattle/internal/app/repository/mongorep"
	"seabattle/internal/entity"
)

//type User interface {
//	GetUser(ctx context.Context, tgID string) (entity.User, error)
//	SetFieldQueryId(ctx context.Context, sessionID, tgId, queryID string, my bool) error
//	SetPoint(ctx context.Context, sessionID, idChatKey string, point entity.Point) error
//	GetPoint(ctx context.Context, sessionID, idChatKey string) (entity.Point, error)
//}

type User struct {
	collection *mongo.Collection
}

func NewUser(collection *mongo.Collection) irep.User {
	return &User{collection: collection}
}
func (db User) SetUser(ctx context.Context, sessionID string, user entity.User) error {
	filter := bson.D{{"session_id", sessionID}}
	_, err := db.collection.UpdateOne(ctx, filter, bson.D{{"$push", bson.D{
		{
			"users", user,
		},
	},
	}})
	if err != nil {
		return err
	}
	return nil

}
func (db User) SetFieldQueryId(ctx context.Context, sessionID, tgId, queryID string, my bool) error {

	var queryField string

	switch my {
	case true:
		queryField = entity.MyFieldQueryId
	case false:
		queryField = entity.EnemyFieldQueryId
	}

	update := bson.D{
		{"$set", bson.D{
			{"users.$[elem]." + queryField, queryID},
		}},
	}
	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.tg_id": tgId},
		},
	})
	filter := bson.D{{"session_id", sessionID}, {"users.tg_id", tgId}}
	_, err := db.collection.UpdateOne(ctx, filter, update, arrayFilters)
	if err != nil {
		return err
	}
	return nil
}

func (db User) GetUser(ctx context.Context, tgID string) (entity.User, error) {
	// Агрегационный конвейер
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"users.tg_id", tgID}}}},
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

	var result struct {
		Users []entity.User `bson:"users"`
	}

	cursor, err := db.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return entity.User{}, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return entity.User{}, err
		}
		// Проверка, что внутри массива users есть хотя бы один элемент
		if len(result.Users) > 0 {
			return result.Users[0], nil
		}
	}

	return entity.User{}, errors.New("user not found")
}

func (db User) SetPoint(ctx context.Context, sessionID, idChatKey string, p entity.Point) error {
	filter := bson.D{{"session_id", sessionID}, {"users.tg_id", idChatKey}}
	update := bson.M{
		"$set": bson.M{
			"users.$[elem].curr_x": p.X,
			"users.$[elem].curr_y": p.Y,
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
func (db User) GetPoint(ctx context.Context, sessionID, tgID string) (entity.Point, error) {
	// Агрегационный конвейер
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"session_id", sessionID}, {"users.tg_id", tgID}}}},
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

	var result struct {
		Users []entity.User `bson:"users"`
	}

	cursor, err := db.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return entity.Point{}, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return entity.Point{}, err
		}
		if len(result.Users) > 0 {
			return entity.Point{X: result.Users[0].CurrX, Y: result.Users[0].CurrY}, nil
		}
	}

	return entity.Point{}, errors.New("user not found")
}
