package app

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/a1k24/short-url/configs"
	"github.com/a1k24/short-url/internal/pkg"
)

func GenerateNextSequence(name string) (int64, error) {
	filter := bson.D{{"name", name}}
	update := bson.D{
		{"$inc", bson.D{
			{"counter", 1},
		}},
	}
	var result CounterSequence
	updateOptions := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	err := pkg.GetCollection(configs.CounterCollectionName).FindOneAndUpdate(context.TODO(), filter, update, updateOptions).Decode(&result)
	if err != nil {
		return 0, err
	}
	return result.Counter, err
}

func DropSequence(name string) (*mongo.DeleteResult, error) {
	filter := bson.D{{"name", name}}
	return pkg.GetCollection(configs.CounterCollectionName).DeleteOne(context.TODO(), filter)
}
