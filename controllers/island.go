package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DATABASE INSTANCE
var collection *mongo.Collection

// IslandCollection : Set global reference to the islands collection.
func IslandCollection(c *mongo.Database) {
	collection = c.Collection("islands")
}

// GetIslands : Get islands from the database (default 10). Currently no order.
func GetIslands(w http.ResponseWriter, r *http.Request) {
	// we created Book array
	// var islands []models.Island

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	cur, _ := collection.Aggregate(ctx, mongo.Pipeline{{{"$lookup", bson.D{{"from", "owners"}, {"localField", "owner"}, {"foreignField", "_id"}, {"as", "owner_details"}}}}}, opts)

	// var results []models.Island
	// if err = cur.All(context.TODO(), &results); err != nil {
	// 	log.Fatal(err)
	// }
	defer cur.Close(ctx)

	// TODO: Figure out how to structure the Owner data from the lookup into an Owner instance...
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// TODO: Get associated owner
		fmt.Print(result)

		// islands = append(islands, result)
	}

	// if err := cur.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	// w.WriteHeader(http.StatusOK)

	// js, err := json.Marshal(islands)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.Write(js)
}
