package configs

import (
	"flag"
	"fmt"
	"time"

	"github.com/allegro/bigcache"
)

var baseUrl string
var DBName string = "globaldb"
var username string
var password string
var dbHost string = "test-cluster.f8tgw.mongodb.net"
var MongoUrl string
var CacheConfig bigcache.Config

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
	CacheConfig = bigcache.Config{
		Shards:             256,
		LifeWindow:         120 * time.Minute,
		CleanWindow:        10 * time.Minute,
		MaxEntriesInWindow: 1000 * 120 * 60,
		MaxEntrySize:       500,
		Verbose:            true,
		HardMaxCacheSize:   256,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}
}

func GetMongoUrl() string {
	return MongoUrl
}

func GetBaseUrl() string {
	return baseUrl
}
