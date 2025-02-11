// package utils

// import (
// 	"context"
// 	"log"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )


// func IsEmailExists(email string, collection *mongo.Collection) (bool, error) {
// 	var result bson.M

// 	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&result)

// 	if err == mongo.ErrNoDocuments {
// 		return false, nil
// 	} else if err != nil {
// 		log.Println("Error checking email:", err)
// 		return false, err

// 	}
// 	return true, nil
// }


package utils