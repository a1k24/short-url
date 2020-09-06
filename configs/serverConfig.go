package configs

import (
	"flag"
	"fmt"
)

var baseUrl string
var DBName string = "globaldb"
var username string
var password string
var dbHost string = "test-cluster.f8tgw.mongodb.net"
var MongoUrl string

const (
	UrlCollectionName     = "shortLink"
	CounterCollectionName = "counterSequence"
)

func Init() {
	flag.StringVar(&username, "username", "username", "user name for mongodb")
	flag.StringVar(&password, "password", "password", "password for mongodb")
	flag.StringVar(&DBName, "dbname", "globaldb", "database name for mongodb")
	flag.StringVar(&baseUrl, "baseurl", "0.0.0.0:10000", "base url for server")
	flag.StringVar(&dbHost, "dbhost", "test-cluster.f8tgw.mongodb.net", "host for mongo server")
	flag.Parse()
	MongoUrl = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", username, password, dbHost, DBName)
}

func GetMongoUrl() string {
	return MongoUrl
}

func GetBaseUrl() string {
	return baseUrl
}
