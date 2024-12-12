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

func PlaceOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var order models.Order
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.OrderID = primitive.NewObjectID()
		order.OrderedAt = time.Now()

		_, err := OrderCollection.InsertOne(ctx, order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not place order"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully"})
	}
}

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orders []models.Order
		cursor, err := OrderCollection.Find(ctx, bson.M{"user_id": userID})
		if err != nil {
			log.Println("Error fetching orders:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching orders"})
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &orders); err != nil {
			log.Println("Error decoding orders:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding orders"})
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}
