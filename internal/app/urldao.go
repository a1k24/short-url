package app

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/a1k24/short-url/configs"
	"github.com/a1k24/short-url/internal/pkg"
)

func SaveUrlToDB(info *UrlInfo) (*mongo.InsertOneResult, error) {
	return pkg.GetCollection(configs.UrlCollectionName).InsertOne(context.TODO(), info)
}

func FindUrlByLinkHash(linkHash string) (*UrlInfo, error) {
	filter := bson.D{{"linkhash", linkHash}}
	var result UrlInfo
	findResult := pkg.GetCollection(configs.UrlCollectionName).FindOne(context.TODO(), filter)
	if findResult.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	err := findResult.Decode(&result)
	if nil != err {
		return nil, err
	}
	return &result, nil
}

func RemoveUrlFromDB(linkHash string) (*mongo.DeleteResult, error) {
	filter := bson.D{{"linkhash", linkHash}}
	return pkg.GetCollection(configs.UrlCollectionName).DeleteOne(context.TODO(), filter)
}

func FindUrlByUrlMd5(urlMd5 string) (*UrlInfo, error) {
	filter := bson.D{{"urlmd5", urlMd5}}
	var result UrlInfo
	findResult := pkg.GetCollection(configs.UrlCollectionName).FindOne(context.TODO(), filter)
	if findResult.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	err := findResult.Decode(&result)
	if nil != err {
		return nil, err
	}
	return &result, nil
}

func IncrementClickCount(linkHash string) {
	filter := bson.D{{"linkhash", linkHash}}
	update := bson.D{
		{"$inc", bson.D{
			{"clickcount", 1},
		}},
	}
	_, err := pkg.GetCollection(configs.UrlCollectionName).UpdateOne(context.TODO(), filter, update)
	if nil != err {
		log.Println("Failed to increment click count: ", err)
	}
}
