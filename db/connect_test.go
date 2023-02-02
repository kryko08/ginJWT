package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

// Test retrieving a document from test collection
func TestConnectDB(t *testing.T) {
	// run function
	client := ConnectDB()
	collection := client.Database("VirtualStorage").Collection("test_collection")

	filter := bson.D{}
	count, err := collection.CountDocuments(context.TODO(), filter)

	if count != 1 || err != nil {
		t.Fatalf("Connection failed")
	}
}
