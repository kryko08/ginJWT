package db

import (
	"GoProject/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func ConnectDB() *mongo.Client {
	mongoUri := utils.GetOsEnv("MONGO_URI")

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(mongoUri).
		SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ping := client.Ping(ctx, readpref.Primary())
	if ping != nil {
		log.Fatal("Error pinging the DB server")
	}

	fmt.Print("Database successfully connected")
	return client
}

// Database client instance
var Database *mongo.Client = ConnectDB()

func GetCollection(collName string, client *mongo.Client) *mongo.Collection {
	collection := client.Database("VirtualStorage").Collection(collName)
	return collection
}
