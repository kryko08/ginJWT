package db

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

func ConnectDB() *mongo.Client {
	// load .env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	mongoUri := os.Getenv("MONGO_URI")

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
