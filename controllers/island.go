package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BrandenSoropia/dodo-getaways-backend/helpers"
	"github.com/BrandenSoropia/dodo-getaways-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DATABASE INSTANCE
var collection *mongo.Collection

// IslandCollection : Set global reference to the islands collection.
func IslandCollection(c *mongo.Database) {
	collection = c.Collection("islands")
}

// Common aggregate steps to populate the Island struct
var lookupStage bson.D = bson.D{{Key: "$lookup", Value: bson.D{
	{Key: "from", Value: "owners"},
	{Key: "localField", Value: "owner_id"},
	{Key: "foreignField", Value: "_id"},
	{Key: "as", Value: "owner_details"},
}}}

var projectStage bson.D = bson.D{{Key: "$project", Value: bson.D{
	{Key: "owner_id", Value: 1},
	{Key: "name", Value: 1},
	{Key: "hemisphere", Value: 1},
	{Key: "description", Value: 1},
	{Key: "is_draft", Value: 1},
	{Key: "images", Value: 1},
	{Key: "owner_details", Value: bson.D{
		{Key: "$arrayElemAt", Value: bson.A{"$owner_details", 0}},
	}},
}}}

type getIslandRequestBody struct {
	IslandID   string `json:"island_id,omitempty" bson:"island_id,omitempty"`
	IslandName string `json:"island_name,omitempty" bson:"island_name,omitempty"`
}

// TODO: Handle no results found by returning empty object
// GetIsland : Return one island matching given island ID or island name from the request.
func GetIsland(w http.ResponseWriter, r *http.Request) {
	var p getIslandRequestBody
	err := helpers.DecodeJSONBody(w, r, &p)

	filter := bson.D{}

	if p.IslandName != "" {
		filter = append(filter, bson.E{Key: "name", Value: p.IslandName})
	}

	if p.IslandID != "" {
		fmt.Print(p.IslandID)
		islandID, _ := primitive.ObjectIDFromHex(p.IslandID)
		filter = append(filter, bson.E{Key: "_id", Value: islandID})
	}

	var island models.Island

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	opts := options.Aggregate().SetMaxTime(2 * time.Second)

	matchStage := bson.D{{Key: "$match", Value: filter}}

	cur, _ := collection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		projectStage,
	}, opts)

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		err := cur.Decode(&island)

		if err != nil {
			log.Fatal(err)
		}
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)

	js, err := json.Marshal(island)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// TODO: Handle no results found by returning empty array
// GetIslands : Get islands from the database (default 10). Currently no order.
func GetIslands(w http.ResponseWriter, r *http.Request) {
	var islands []models.Island

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	cur, _ := collection.Aggregate(ctx, mongo.Pipeline{
		lookupStage,
		projectStage,
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

type setIslandRequestBody struct {
	ID          string   `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerID     string   `json:"owner_id" bson:"owner_id"`
	Name        string   `json:"name" bson:"name"`
	Hemisphere  string   `json:"hemisphere" bson:"hemisphere"`
	Description string   `json:"description" bson:"description"`
	Images      []string `json:"images" bson:"images"`
	IsDraft     bool     `json:"is_draft" bson:"is_draft"`
}

// SetIsland : Insert a new island into the islands collection, or if given an existing island id, update it.
func SetIsland(w http.ResponseWriter, r *http.Request) {
	var body setIslandRequestBody
	err := helpers.DecodeJSONBody(w, r, &body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ownerObjectID, err := primitive.ObjectIDFromHex(body.OwnerID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$exists", Value: false}}}}

	newIsland := bson.D{
		{Key: "owner_id", Value: ownerObjectID},
		{Key: "name", Value: body.Name},
		{Key: "hemisphere", Value: body.Hemisphere},
		{Key: "description", Value: body.Description},
		{Key: "is_draft", Value: body.IsDraft},
		{Key: "images", Value: body.Images},
	}

	// If givin an island id, this controller will now update the matching island
	if body.ID != "" {
		islandObjectID, err := primitive.ObjectIDFromHex(body.ID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newIsland = append(newIsland, bson.E{Key: "_id", Value: islandObjectID})
		filter = bson.D{{Key: "_id", Value: islandObjectID}}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var updatedDocument models.Island

	update := bson.D{{
		Key: "$set", Value: newIsland,
	}}

	// TODO: Check that owner exists first and does not have an island already. Do not allow if owner DNE, or if they have an island but gave a non-matching island id

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)

	// TODO: Clean up handling update vs new island...
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusOK)
			js, err := json.Marshal(updatedDocument)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)

		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
		}
	} else {

		w.WriteHeader(http.StatusOK)

		js, err := json.Marshal(updatedDocument)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
