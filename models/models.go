package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Island : Struct for island data from DB.
type Island struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	OwnerID      primitive.ObjectID `json:"owner_id" bson:"owner_id"`
	Name         string             `json:"name" bson:"name"`
	Hemisphere   string             `json:"hemisphere" bson:"hemisphere"`
	Description  string             `json:"description" bson:"description"`
	IsDraft      bool               `json:"is_draft" bson:"isDraft"`
	Images       []string           `json:"images" bson:"images"`
	OwnerDetails Owner              `json:"owner_details,omitempty" bson:"owner_details,omitempty"`
}

// Owner : Struct for owner data from DB.
type Owner struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	IslandID        primitive.ObjectID `json:"island_id" bson:"island_id"`
	Username        string             `json:"username" bson:"username"`
	DiscordUsername string             `json:"discord_username" bson:"discord_username"`
	Description     string             `json:"description" bson:"description"`
}
