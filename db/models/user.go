package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user document in the MongoDB collection.
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // MongoDB Object ID
	Username  string             `bson:"username" json:"username"`          // Username of the user
	Email     string             `bson:"email" json:"email"`                // Email address
	Password  string             `bson:"password" json:"password"`          // Hashed password
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`      // Creation timestamp
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`      // Update timestamp
}
