package main

import (
	"log"

	"github.com/a1k24/short-url/configs"
	"github.com/a1k24/short-url/internal/app"
	"github.com/a1k24/short-url/internal/pkg"
)

func main() {
	configs.Init()
	log.Println("Server started.")
	_, cancel := pkg.CreateConnection() // ensure mongo client is created at start
	log.Println("Mongo connection established.")
	defer cancel()
	app.HandleRequests()
}
