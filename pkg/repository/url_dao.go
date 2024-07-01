package repository

import (
	"context"
	"shortener/pkg"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DAO - data acsess object - main DB logic
// create new mongo DB collection
type UrlDAO struct {
	db *mongo.Collection
}

func (dao *UrlDAO) Insert(ctx context.Context, shortURL *pkg.ShortURL) error {
	_, err := dao.db.InsertOne(ctx, shortURL)
	return err
}

func (dao *UrlDAO) FindByID(ctx context.Context, id string) (*pkg.ShortURL, error) {
	filter := bson.D{{"_id", id}}
	var shortURL pkg.ShortURL
	err := dao.db.FindOne(ctx, filter).Decode(&shortURL)
	switch {
	case err == nil:
		return &shortURL, nil
	case err == mongo.ErrNoDocuments:
		return nil, err
	default:
		return nil, err
	}
}
