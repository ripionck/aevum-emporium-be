package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID      primitive.ObjectID `bson:"_id" json:"user_id"`
	FirstName   string             `bson:"first_name" json:"first_name"`
	LastName    string             `bson:"last_name" json:"last_name"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"password"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number"`
	Address     []Address          `bson:"address" json:"address"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	Role        string             `bson:"role" json:"role"` // "admin" or "customer"
}
