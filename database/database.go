package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() {

	MongoDB := os.Getenv("MONGODB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongodbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDB))

	if err != nil {
		log.Fatalln("error while connecting to database")
	}

	log.Println("Connected to Database Successfully")

	Client = mongodbClient
}

func OpenCollection(collectionName string) *mongo.Collection {
	collection := Client.Database(os.Getenv("auth")).Collection(collectionName)

	return collection
}
