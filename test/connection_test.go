package test

import (
	"testing"

	"github.com/a1k24/short-url/internal/pkg"
)

func TestCreateConnection(t *testing.T) {
	connection, cancelFunc := pkg.CreateConnection()
	if nil == connection {
		t.Error("Failed to create connection.")
	}
	if nil == cancelFunc {
		t.Error("Invalid cancel Func")
	}
	defer cancelFunc()
}

func TestGetMongoClient(t *testing.T) {
	client := pkg.GetMongoClient()
	if nil == client {
		t.Error("Failed to get mongo client.")
	}
	mongoClient := pkg.GetMongoClient()

	if mongoClient != client {
		t.Error("Got different client instance.")
	}

}
