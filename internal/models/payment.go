package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	PaymentID primitive.ObjectID `bson:"_id" json:"payment_id"`
	Method    string             `bson:"method" json:"method"` //  "Credit Card", "PayPal")
	Status    string             `bson:"status" json:"status"` //  "Pending", "Completed")
	Amount    float64            `bson:"amount" json:"amount"`
	PaidAt    *time.Time         `bson:"paid_at,omitempty" json:"paid_at,omitempty"`
}
