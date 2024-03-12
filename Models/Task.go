package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id                primitive.ObjectID `bson:"_id,omitempty"` //MongoDb document ID
	UserId            primitive.ObjectID `bson:"user_id"`
	Title             string             `bson:"title"`
	Description       string             `bson:"description"`
	Status            TaskStatus         `bson:"status"`
	DateTimeCreated   primitive.DateTime `bson:"date_time_created"`
	DateTimeCompleted primitive.DateTime `bson:"date_time_completed"`
}
