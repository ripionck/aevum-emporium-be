package controllers

import (
	"aevum-emporium-be/internal/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure user is authenticated
		userID := c.GetString("uid")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Convert userID to ObjectID
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}

		// Bind the new address
		var newAddress models.Address
		newAddress.AddressID = primitive.NewObjectID()
		if err := c.BindJSON(&newAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address format"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Ensure the user doesn't have more than 2 addresses
		countFilter := bson.D{{Key: "_id", Value: userObjectID}}
		countPipeline := mongo.Pipeline{
			{{Key: "$match", Value: countFilter}},
			{{Key: "$unwind", Value: "$address"}},
			{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			}}},
		}

		cursor, err := UserCollection.Aggregate(ctx, countPipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check address count"})
			return
		}

		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing address count"})
			return
		}

		var addressCount int32 = 0
		if len(results) > 0 {
			addressCount = results[0]["count"].(int32)
		}

		// Allow a maximum of 2 addresses
		if addressCount >= 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum of 2 addresses allowed"})
			return
		}

		// Add the new address to the user's document
		filter := bson.D{{Key: "_id", Value: userObjectID}}
		update := bson.D{{Key: "$push", Value: bson.D{{Key: "address", Value: newAddress}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add address"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Address added successfully"})
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure user is authenticated
		userID := c.GetString("uid")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Convert userID to ObjectID
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}

		// Bind the updated address
		var updatedAddress models.Address
		if err := c.BindJSON(&updatedAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address format"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Check if the address being edited is the home address
		filter := bson.D{{Key: "_id", Value: userObjectID}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "address.0.street", Value: updatedAddress.Street},
			{Key: "address.0.city", Value: updatedAddress.City},
			{Key: "address.0.state", Value: updatedAddress.State},
			{Key: "address.0.country", Value: updatedAddress.Country},
			{Key: "address.0.zip_code", Value: updatedAddress.ZipCode},
			{Key: "address.0.is_default", Value: updatedAddress.IsDefault},
		}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Home address updated successfully"})
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure user is authenticated
		userID := c.GetString("uid")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Convert userID to ObjectID
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}

		// Bind the updated address
		var updatedAddress models.Address
		if err := c.BindJSON(&updatedAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address format"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Check if the address being edited is the work address
		filter := bson.D{{Key: "_id", Value: userObjectID}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "address.1.street", Value: updatedAddress.Street},
			{Key: "address.1.city", Value: updatedAddress.City},
			{Key: "address.1.state", Value: updatedAddress.State},
			{Key: "address.1.country", Value: updatedAddress.Country},
			{Key: "address.1.zip_code", Value: updatedAddress.ZipCode},
			{Key: "address.1.is_default", Value: updatedAddress.IsDefault},
		}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Work address updated successfully"})
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure user is authenticated
		userID := c.GetString("uid")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Convert userID to ObjectID
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Remove all addresses for the user
		filter := bson.D{{Key: "_id", Value: userObjectID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "address", Value: []models.Address{}}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete addresses"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Addresses deleted successfully"})
	}
}
