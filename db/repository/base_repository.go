package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseRepository struct {
	Collection *mongo.Collection
}

// Create inserts a new document into the collection.
func (r *BaseRepository) Create(ctx context.Context, document interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.Collection.InsertOne(ctx, document)
	return err
}

// Find retrieves documents matching the given filter.
func (r *BaseRepository) Find(ctx context.Context, filter interface{}, results interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, results)
}

// FindOne retrieves a single document matching the filter.
func (r *BaseRepository) FindOne(ctx context.Context, filter interface{}, result interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.Collection.FindOne(ctx, filter).Decode(result)
}

// Update updates documents matching the filter.
func (r *BaseRepository) Update(ctx context.Context, filter interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.Collection.UpdateMany(ctx, filter, update)
	return err
}

// Delete removes documents matching the filter.
func (r *BaseRepository) Delete(ctx context.Context, filter interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.Collection.DeleteMany(ctx, filter)
	return err
}

// GetAll retrieves all documents from the collection and decodes them into the provided results.
func (r *BaseRepository) FindAll(ctx context.Context, results interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, results)
}
