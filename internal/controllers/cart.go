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

var CartCollection *mongo.Collection = datasource.CartData(datasource.Client)

// AddToCart adds a product to the user's cart
func AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var cartItem models.CartItem
		if err := c.BindJSON(&cartItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID := c.GetString("user_id") // Assume user ID is extracted from token middleware
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		userObjID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var cart models.Cart
		err = CartCollection.FindOne(ctx, bson.M{"user_id": userObjID}).Decode(&cart)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// Create a new cart if none exists
				cart = models.Cart{
					CartID:    primitive.NewObjectID(),
					UserID:    userObjID,
					Items:     []models.CartItem{cartItem},
					Total:     cartItem.Price * float64(cartItem.Quantity),
					CreatedAt: time.Now(),
				}
				_, err := CartCollection.InsertOne(ctx, cart)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating cart"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "Item added to cart", "cart": cart})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving cart"})
				return
			}
		}

		// Check if the item already exists in the cart
		itemExists := false
		for i, item := range cart.Items {
			if item.ProductID == cartItem.ProductID {
				cart.Items[i].Quantity += cartItem.Quantity
				cart.Items[i].Price = cartItem.Price
				itemExists = true
				break
			}
		}

		if !itemExists {
			cart.Items = append(cart.Items, cartItem)
		}

		// Update the total price
		cart.Total = 0
		for _, item := range cart.Items {
			cart.Total += item.Price * float64(item.Quantity)
		}

		_, err = CartCollection.UpdateOne(ctx, bson.M{"_id": cart.CartID}, bson.M{"$set": bson.M{"items": cart.Items, "total": cart.Total}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating cart"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Item added to cart", "cart": cart})
	}
}

// ViewCart fetches the cart items for a user
func ViewCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var cart models.Cart
		err := CartCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
		if err != nil {
			log.Println("Error fetching cart:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching cart"})
			return
		}

		c.JSON(http.StatusOK, cart)
	}
}

// RemoveFromCart removes a specific product from the user's cart
func RemoveFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID := c.GetString("user_id") // Assume user ID is extracted from token middleware
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		userObjID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		productID := c.Param("product_id")
		productObjID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var cart models.Cart
		err = CartCollection.FindOne(ctx, bson.M{"user_id": userObjID}).Decode(&cart)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving cart"})
			return
		}

		// Remove the item from the cart
		itemRemoved := false
		for i, item := range cart.Items {
			if item.ProductID == productObjID {
				// Remove item from slice
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
				itemRemoved = true
				break
			}
		}

		if !itemRemoved {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found in cart"})
			return
		}

		// Recalculate the total price
		cart.Total = 0
		for _, item := range cart.Items {
			cart.Total += item.Price * float64(item.Quantity)
		}

		_, err = CartCollection.UpdateOne(ctx, bson.M{"_id": cart.CartID}, bson.M{"$set": bson.M{"items": cart.Items, "total": cart.Total}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating cart"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart", "cart": cart})
	}
}

// ClearCart removes all items from the user's cart
func ClearCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID := c.GetString("user_id") // Assume user ID is extracted from token middleware
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		userObjID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var cart models.Cart
		err = CartCollection.FindOne(ctx, bson.M{"user_id": userObjID}).Decode(&cart)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving cart"})
			return
		}

		// Clear the cart (empty items and reset total)
		cart.Items = []models.CartItem{}
		cart.Total = 0

		_, err = CartCollection.UpdateOne(ctx, bson.M{"_id": cart.CartID}, bson.M{"$set": bson.M{"items": cart.Items, "total": cart.Total}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error clearing cart"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cart cleared", "cart": cart})
	}
}
