package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	OrderID       primitive.ObjectID `bson:"_id" json:"order_id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items         []OrderItem        `bson:"items" json:"items"`
	TotalPrice    float64            `bson:"total_price" json:"total_price"`
	Discount      *float64           `bson:"discount" json:"discount"`
	OrderedAt     time.Time          `bson:"ordered_at" json:"ordered_at"`
	PaymentMethod Payment            `bson:"payment_method" json:"payment_method"`
	Status        string             `bson:"status" json:"status"` //"Processing", "Shipped", "Delivered"
	TransactionID string             `bson:"transaction_id" json:"transaction_id"`
	PaymentStatus string             `bson:"payment_status" json:"payment_status"` //"Pending", "Succeeded", "Failed"
}

type OrderItem struct {
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	Name      string             `bson:"name" json:"name"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	Price     float64            `bson:"price" json:"price"`
}
