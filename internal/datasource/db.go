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

// ConnectDB establishes a connection to MongoDB and returns the client.
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

// Global variable to hold the MongoDB client
var Client *mongo.Client = ConnectDB()

func getCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("aevum-emporium").Collection(collectionName)
}

func UserData(client *mongo.Client) *mongo.Collection {
	return getCollection(client, "User")
}

func ProductData(client *mongo.Client) *mongo.Collection {
	return getCollection(client, "Product")
}

func OrderData(client *mongo.Client) *mongo.Collection {
	return getCollection(client, "Order")
}

func CartData(client *mongo.Client) *mongo.Collection {
	return getCollection(client, "Cart")
}

func ReviewData(client *mongo.Client) *mongo.Collection {
	return getCollection(client, "Review")
}

func WishlistData(client *mongo.Client) *mongo.Collection {
	return getCollection(client, "Wishlist")
}
