package models

type User struct {
	ID       string `bson:"_id,omitempty"`
	Username string `bson:"username" binding:"required"`
	Password string `bson:"password" binding:"required"`
}
