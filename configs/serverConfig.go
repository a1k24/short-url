package configs

import (
	"flag"
	"fmt"
)

var baseUrl string
var DBName string = "globaldb"
var username string
var password string
var MongoUrl string

const (
	UrlCollectionName     = "shortLink"
	CounterCollectionName = "counterSequence"
)

func Init() {
	flag.StringVar(&username, "username", "username", "user name for mongodb")
	flag.StringVar(&password, "password", "password", "password for mongodb")
	flag.StringVar(&DBName, "dbname", "globaldb", "database name for mongodb")
	flag.StringVar(&baseUrl, "base_url", "localhost:10000", "base url for server")
	flag.Parse()
	MongoUrl = fmt.Sprintf("mongodb+srv://%s:%s@test-cluster.f8tgw.mongodb.net/%s?retryWrites=true&w=majority", username, password, DBName)
}

func GetMongoUrl() string {
	return MongoUrl
}

func GetBaseUrl() string {
	return baseUrl
}
