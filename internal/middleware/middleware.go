package middleware

import (
	"net/http"
	"strings"

	token "aevum-emporium-be/internal/token"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a Gin middleware that checks for a valid Authorization token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the token from the "Authorization" header
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			// No Authorization header provided
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the token starts with "Bearer " (standard format)
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Extract the token from the "Bearer <token>" format
		ClientToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			// If there's an error validating the token, respond with Unauthorized status
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		// Store the claims in the context
		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)

		// Continue to the next handler
		c.Next()
	}
}
