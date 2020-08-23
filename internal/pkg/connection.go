package pkg

import (
	"context"
	"github.com/a1k24/short-url/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var client *mongo.Client = nil
var cancel context.CancelFunc = nil

func CreateConnection() (*mongo.Client, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		configs.MongoUrl,
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
