package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ProductID     primitive.ObjectID `bson:"_id" json:"product_id"`
	Name          string             `bson:"name" json:"name"`
	Description   string             `bson:"description" json:"description"`
	Price         float64            `bson:"price" json:"price"`
	StockQuantity int                `bson:"stock_quantity" json:"stock_quantity"`
	Category      string             `bson:"category" json:"category"`
	Images        []string           `bson:"images" json:"images"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	Discount      *float64           `bson:"discount" json:"discount"`
}
