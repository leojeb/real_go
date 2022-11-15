package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	var coll *mongo.Collection
	var firstID, secondID primitive.ObjectID

	// Update the "email" field for two users.
	// For each update, specify the Upsert option to insert a new document if a
	// document matching the filter isn't found.
	// Set the Ordered option to false to allow both operations to happen even
	// if one of them errors.
	firstUpdate := bson.D{
		{"$set", bson.D{
			{"email", "firstEmail@example.com"},
		}},
	}
	secondUpdate := bson.D{
		{"$set", bson.D{
			{"email", "secondEmail@example.com"},
		}},
	}
	models := []mongo.WriteModel{
		mongo.NewUpdateOneModel().SetFilter(bson.D{{"_id", firstID}}).
			SetUpdate(firstUpdate).SetUpsert(true),
		mongo.NewUpdateOneModel().SetFilter(bson.D{{"_id", secondID}}).
			SetUpdate(secondUpdate).SetUpsert(true),
	}
	opts := options.BulkWrite().SetOrdered(false)
	res, err := coll.BulkWrite(context.TODO(), models, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"inserted %v and deleted %v documents\n",
		res.InsertedCount,
		res.DeletedCount)
}
