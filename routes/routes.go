package routes

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	OrderID   string             `bson:"order_id" json:"order_id"`
	RiderID   string             `bson:"rider_id,omitempty" json:"rider_id,omitempty"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}
