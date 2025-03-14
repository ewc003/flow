package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Breakdown
type Breakdown struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // MongoDB Object ID
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`            // Reference to the user who owns this breakdown
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
