package utils

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitDB() error {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return Err("MONGO_URI not set")
	}

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	return nil
}

func GetCollection(name string) *mongo.Collection {
	return Client.Database("urlshortener").Collection(name)
}

func Err(msg string) error {
	log.Println("❌", msg)
	return &customErr{msg}
}

type customErr struct {
	msg string
}

func (e *customErr) Error() string {
	return e.msg
}
