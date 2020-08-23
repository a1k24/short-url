package app

import (
	"context"

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
