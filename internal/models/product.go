package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id          primitive.ObjectID `json:"id"          bson:"_id"`
	Name        string             `json:"name"        bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float32            `json:"price"       bson:"price"`
	CreatedAt   int64              `json:"-"           bson:"CreatedAt"`
	UpdatedAt   int64              `json:"-"           bson:"UpdatedAt"`
	DeletedAt   int64              `json:"-"           bson:"DeletedAt"`
}
