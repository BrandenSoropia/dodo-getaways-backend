package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/BrandenSoropia/dodo-getaways-backend/models"
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
	var islands []models.Island

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	cur, _ := collection.Aggregate(ctx, mongo.Pipeline{
		{{"$lookup", bson.D{
			{"from", "owners"},
			{"localField", "owner_id"},
			{"foreignField", "_id"},
			{"as", "owner_details"},
		}}},
		{{"$project", bson.D{
			{"owner_id", 1},
			{"name", 1},
			{"hemisphere", 1},
			{"description", 1},
			{"is_draft", 1},
			{"images", 1},
			{"owner_details", bson.D{
				{"$arrayElemAt", bson.A{"$owner_details", 0}},
			}},
		}}},
	}, opts)

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var island models.Island
		err := cur.Decode(&island)

		if err != nil {
			log.Fatal(err)
		}

		islands = append(islands, island)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)

	js, err := json.Marshal(islands)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
