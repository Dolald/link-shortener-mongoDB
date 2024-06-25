package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// DAO - data acsess object - main DB logic
type UrlDAO struct {
	db *mongo.Collection
}

func NewUrlDAO(ctx context.Context, client *mongo.Client) (*UrlDAO, error) {
	dao := &UrlDAO{
		db: client.Database("core").Collection("shortUrls"),
	}
	return dao, nil
}
