package configs

import "fmt"

var BaseUrl = "localhost:10000"
var DBName = "<enter dbname here>"
var username = "<enter username here>"
var password = "<enter password here>"
var MongoUrl = fmt.Sprintf("mongodb+srv://%s:%s@test-cluster.f8tgw.mongodb.net/%s?retryWrites=true&w=majority", username, password, DBName)
