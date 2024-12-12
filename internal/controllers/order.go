package controllers

import (
	"aevum-emporium-be/internal/datasource"
	"aevum-emporium-be/internal/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var OrderCollection *mongo.Collection = datasource.OrderData(datasource.Client)

// PlaceOrder adds a new order to the database for the authenticated user
func PlaceOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure user is authenticated
		userID := c.GetString("uid")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Convert the userID to ObjectID
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var order models.Order
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Ensure the order has the user ID
		order.UserID = userObjectID
		order.OrderID = primitive.NewObjectID()
		order.OrderedAt = time.Now()

		// Insert the order into the database
		_, err = OrderCollection.InsertOne(ctx, order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not place order"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order_id": order.OrderID})
	}
}

// GetOrders retrieves all orders for the authenticated user
func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure user is authenticated
		userID := c.GetString("uid")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Convert the userID to ObjectID
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Find orders belonging to the authenticated user
		var orders []models.Order
		cursor, err := OrderCollection.Find(ctx, bson.M{"user_id": userObjectID})
		if err != nil {
			log.Println("Error fetching orders:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching orders"})
			return
		}
		defer cursor.Close(ctx)

		// Decode the orders
		if err := cursor.All(ctx, &orders); err != nil {
			log.Println("Error decoding orders:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding orders"})
			return
		}

		// If no orders found, return an appropriate message
		if len(orders) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No orders found for this user"})
			return
		}

		// Return orders if found
		c.JSON(http.StatusOK, orders)
	}
}
