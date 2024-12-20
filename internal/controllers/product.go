package controllers

import (
	"aevum-emporium-be/internal/datasource"
	"aevum-emporium-be/internal/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ProductCollection *mongo.Collection = datasource.ProductData(datasource.Client)

func isAdmin(userID string) (bool, error) {
	// Convert userID string to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user ID format")
	}

	fmt.Println("Checking user role for UID:", userID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	// Query based on ObjectID
	err = UserCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil // No user found, not an admin
		}
		return false, err // Other errors (DB issues)
	}

	fmt.Println("User role found:", user.Role)

	// Check if the user role is admin
	return user.Role == "admin", nil
}

func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure user is authenticated and has admin role
		userID := c.GetString("uid")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		isAdmin, err := isAdmin(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user role"})
			return
		}
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add products"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var product models.Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set product ID and timestamps
		product.ProductID = primitive.NewObjectID()
		product.CreatedAt = time.Now()
		product.UpdatedAt = time.Now()

		_, err = ProductCollection.InsertOne(ctx, product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Product could not be created"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product successfully created"})
	}
}

func GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var products []models.Product
		cursor, err := ProductCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Println("Error fetching products:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &products); err != nil {
			log.Println("Error decoding products:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding products"})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}

func GetProductByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var product models.Product
		err = ProductCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
		if err != nil {
			log.Println("Product not found:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure user is authenticated and has admin role
		userID := c.GetString("uid")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		isAdmin, err := isAdmin(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user role"})
			return
		}
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update products"})
			return
		}

		productID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var product models.Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product.UpdatedAt = time.Now()

		update := bson.M{"$set": bson.M{
			"name":           product.Name,
			"category":       product.Category,
			"description":    product.Description,
			"price":          product.Price,
			"stock_quantity": product.StockQuantity,
			"updated_at":     product.UpdatedAt,
		}}

		_, err = ProductCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			log.Println("Error updating product:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product successfully updated"})
	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure user is authenticated and has admin role
		userID := c.GetString("uid")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		isAdmin, err := isAdmin(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user role"})
			return
		}
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete products"})
			return
		}

		productID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		_, err = ProductCollection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			log.Println("Error deleting product:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product successfully deleted"})
	}
}

func SearchProductByCategoryOrName() gin.HandlerFunc {
	return func(c *gin.Context) {
		category := c.Query("category")
		name := c.Query("name")

		var filter bson.M
		if category != "" && name != "" {
			filter = bson.M{"$and": []bson.M{
				{"category": bson.M{"$regex": category, "$options": "i"}},
				{"name": bson.M{"$regex": name, "$options": "i"}},
			}}
		} else if category != "" {
			filter = bson.M{"category": bson.M{"$regex": category, "$options": "i"}}
		} else if name != "" {
			filter = bson.M{"name": bson.M{"$regex": name, "$options": "i"}}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide either category or name to search"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var products []models.Product
		cursor, err := ProductCollection.Find(ctx, filter)
		if err != nil {
			log.Println("Error fetching products by category or name:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
			return
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &products); err != nil {
			log.Println("Error decoding products:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding products"})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}
