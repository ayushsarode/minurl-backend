package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URL struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Original  string             `bson:"original" json:"original"`
	Short     string             `bson:"short" json:"short"`
	UserID    string             `bson:"userID" json:"userID"`
	Clicks    int                `bson:"clicks" json:"clicks"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
