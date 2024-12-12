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

var ReviewCollection *mongo.Collection = datasource.ReviewData(datasource.Client)

func AddReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var review models.Review
		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		review.ReviewID = primitive.NewObjectID()
		review.CreatedAt = time.Now()

		_, err := ReviewCollection.InsertOne(ctx, review)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding review"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Review added successfully", "review": review})
	}
}

func GetReviewsByProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		productID := c.Param("product_id")
		objID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var reviews []models.Review
		cursor, err := ReviewCollection.Find(ctx, bson.M{"product_id": objID})
		if err != nil {
			log.Println("Error fetching reviews:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching reviews"})
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &reviews); err != nil {
			log.Println("Error decoding reviews:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding reviews"})
			return
		}

		c.JSON(http.StatusOK, reviews)
	}
}

func DeleteReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		reviewID := c.Param("review_id")
		objID, err := primitive.ObjectIDFromHex(reviewID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		_, err = ReviewCollection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			log.Println("Error deleting review:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting review"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
	}
}
