package controllers

import (
	"aevum-emporium-be/internal/datasource"
	"aevum-emporium-be/internal/models"
	generate "aevum-emporium-be/internal/token"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = datasource.UserData(datasource.Client)
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Or Passowrd is Incorerct"
		valid = false
	}
	return valid, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if email already exists
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking for existing user"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}

		// Check if phone number already exists
		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone_number": user.PhoneNumber})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking for existing phone number"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already in use"})
			return
		}

		// Hash the user's password
		hashedPassword := HashPassword(user.Password)
		user.Password = hashedPassword

		// Set timestamps and create user ID
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.UserID = primitive.NewObjectID()

		// Default role to "customer" if not provided
		if user.Role == "" {
			user.Role = "customer"
		}

		// Initialize address list if not provided
		if user.Address == nil {
			user.Address = make([]models.Address, 0)
		}

		// Insert the user into the database
		_, insertErr := UserCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Successfully signed up!"})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var userLoginDetails struct {
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required"`
		}

		var foundUser models.User

		if err := c.BindJSON(&userLoginDetails); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate input
		validationErr := Validate.Struct(userLoginDetails)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Find the user in the database
		err := UserCollection.FindOne(ctx, bson.M{"email": userLoginDetails.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Verify password
		isValidPassword, msg := VerifyPassword(userLoginDetails.Password, foundUser.Password)
		if !isValidPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		// Generate tokens and update them in the database
		token, refreshToken, _ := generate.TokenGenerator(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.UserID.Hex())
		generate.UpdateAllTokens(token, refreshToken, foundUser.UserID.Hex())

		// Send user information back to the client (excluding sensitive fields)
		c.JSON(http.StatusOK, gin.H{
			"user_id":       foundUser.UserID,
			"first_name":    foundUser.FirstName,
			"last_name":     foundUser.LastName,
			"email":         foundUser.Email,
			"phone_number":  foundUser.PhoneNumber,
			"address":       foundUser.Address,
			"role":          foundUser.Role,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}
