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
	// we created Book array
	var islands []models.Island

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// bson.M{},  we passed empty filter. So we want to get all data.
	opts := options.Find().SetLimit(2)
	cur, err := collection.Find(ctx, bson.M{}, opts)

	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result models.Island
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// TODO: Get associated owner

		islands = append(islands, result)
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
