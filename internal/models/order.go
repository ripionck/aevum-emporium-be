package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id"`
	Order_Cart     []Cart             `json:"order_list"  bson:"order_list"`
	Ordered_At     time.Time          `json:"ordered_on"  bson:"ordered_on"`
	Price          int                `json:"total_price" bson:"total_price"`
	Discount       *int               `json:"discount"    bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}
