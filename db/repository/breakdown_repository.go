package repository

import "go.mongodb.org/mongo-driver/mongo"

type BreakdownRepository struct {
	BaseRepository
}

func NewBreakdownRepository(db *mongo.Client) *BreakdownRepository {
	return &BreakdownRepository{
		BaseRepository{
			Collection: db.Database("flow").Collection("breakdowns"),
		},
	}
}
