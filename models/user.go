package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username   string             `bson:"username" json:"username"`
	Email      string             `bson:"email" json:"email" binding:"required,email"`
	Password   string             `bson:"password" json:"password" binding:"required"`
	ProfilePic string             `bson:"profile_pic,omitempty" json:"profile_pic"`
}
