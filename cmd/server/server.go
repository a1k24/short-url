package main

import (
	"github.com/a1k24/short-url/internal/app"
	"log"
)

func main() {
	log.Println("Server started.")
	app.HandleRequests()
}
