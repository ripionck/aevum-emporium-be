package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wishlist struct {
	WishlistID primitive.ObjectID   `bson:"_id" json:"wishlist_id"`
	UserID     primitive.ObjectID   `bson:"user_id" json:"user_id"`
	Products   []primitive.ObjectID `bson:"products" json:"products"`
	CreatedAt  time.Time            `bson:"created_at" json:"created_at"`
}
