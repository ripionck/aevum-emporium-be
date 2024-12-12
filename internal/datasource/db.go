package datasource

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not loaded. Ensure environment variables are set.")
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("Error: MONGO_URI is not set in environment variables or .env file")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB")
	return client
}

var Client *mongo.Client = ConnectDB()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	var userCollection *mongo.Collection = client.Database("aevum-emporium").Collection(collectionName)
	return userCollection

}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var productCollection *mongo.Collection = client.Database("aevum-emporium").Collection(collectionName)
	return productCollection
}
