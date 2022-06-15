package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notifier struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	UserId int                `json:"userId" bson:"user_id"`
}
