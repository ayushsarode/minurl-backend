package models

type URL struct {
	ID       string `bson:"_id,omitempty"`
	Original string `bson:"original"`
	Short    string `bson:"short"`
	UserID   string `bson:"userID"`
}

