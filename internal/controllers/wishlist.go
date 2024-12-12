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

var WishlistCollection *mongo.Collection = datasource.UserData(datasource.Client, "Wishlist")

// AddWishlist adds a product to the user's wishlist
func AddWishlist() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var wishlist models.Wishlist
		if err := c.BindJSON(&wishlist); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Assign a new ObjectID to Wishlist
		wishlist.WishlistID = primitive.NewObjectID()
		wishlist.CreatedAt = time.Now()

		// Insert the wishlist document into the Wishlist collection
		_, err := WishlistCollection.InsertOne(ctx, wishlist)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding to wishlist"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product added to wishlist", "wishlist": wishlist})
	}
}

// GetWishlistByUser fetches all products in a user's wishlist
func GetWishlistByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID := c.Param("user_id")
		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var wishlist models.Wishlist
		err = WishlistCollection.FindOne(ctx, bson.M{"user_id": objID}).Decode(&wishlist)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"message": "No wishlist found for this user"})
				return
			}
			log.Println("Error fetching wishlist:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching wishlist"})
			return
		}

		c.JSON(http.StatusOK, wishlist)
	}
}

// DeleteWishlist removes a product from the user's wishlist
func DeleteWishlist() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		wishlistID := c.Param("wishlist_id")
		objID, err := primitive.ObjectIDFromHex(wishlistID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID"})
			return
		}

		// Remove the wishlist document from the Wishlist collection
		_, err = WishlistCollection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			log.Println("Error deleting from wishlist:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting from wishlist"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product removed from wishlist"})
	}
}
