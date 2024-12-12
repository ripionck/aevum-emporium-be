package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Category     *string            `json:"category"`
	Price        *uint64            `json:"price"`
	Rating       *uint8             `json:"rating"`
	Description  *string            `json:"description"`
	Image        *string            `json:"image"`
}
