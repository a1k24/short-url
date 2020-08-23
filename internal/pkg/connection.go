package pkg

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/a1k24/short-url/configs"
)

var client *mongo.Client = nil
var cancel context.CancelFunc = nil

func CreateConnection() (*mongo.Client, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		configs.GetMongoUrl(),
	))
	if err != nil {
		log.Fatal(err)
	}
	return client, cancel
}

func GetMongoClient() *mongo.Client {
	if client == nil {
		cancelConnection() // cancel previous connection if any
		client, cancel = CreateConnection()
	}
	return client
}

func cancelConnection() {
	if nil != cancel {
		cancel()
	}
}

func getGlobalDB() *mongo.Database {
	return GetMongoClient().Database(configs.GetDBName())
}

func GetCollection(name string) *mongo.Collection {
	return getGlobalDB().Collection(name)
}
