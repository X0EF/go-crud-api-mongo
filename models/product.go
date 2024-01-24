package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id          primitive.ObjectID `json:"id"          bson:"_id"`
	Name        string             `json:"name"        bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float32            `json:"price"       bson:"price"`
	CreatedAt   string             `json:"-"           bson:"CreatedAt"`
	UpdatedAt   string             `json:"-"           bson:"UpdatedAt"`
	DeletedAt   string             `json:"-"           bson:"DeletedAt"`
}
