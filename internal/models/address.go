package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	AddressID primitive.ObjectID `bson:"_id" json:"address_id"`
	Street    string             `bson:"street" json:"street"`
	City      string             `bson:"city" json:"city"`
	State     string             `bson:"state" json:"state"`
	Country   string             `bson:"country" json:"country"`
	ZipCode   string             `bson:"zip_code" json:"zip_code"`
	IsDefault bool               `bson:"is_default" json:"is_default"`
}
