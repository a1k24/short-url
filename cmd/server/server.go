package main

import (
	"github.com/a1k24/short-url/internal/app"
	"github.com/a1k24/short-url/internal/pkg"
	"log"
)

func main() {
	log.Println("Server started.")
	_, cancel := pkg.CreateConnection() // ensure mongo client is created at start
	log.Println("Mongo connection established.")
	defer cancel()
	app.HandleRequests()
}
