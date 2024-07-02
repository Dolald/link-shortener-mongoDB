package repository

import (
	"context"
	"fmt"
	"shortener/pkg"
	domain "shortener/pkg/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Insert(ctx context.Context, shortURL *domain.ShortURL) error
	FindByID(ctx context.Context, id string) (*domain.ShortURL, error)
	Update(ctx context.Context, shortUrl *domain.ShortURL) error
}

// DAO - data acsess object - main DB logic
// create new mongo DB collection
type repository struct {
	db *mongo.Collection
}

func (r *repository) Insert(ctx context.Context, shortURL *domain.ShortURL) error {
	_, err := r.db.InsertOne(ctx, shortURL)

	return err
}

func (r *repository) FindByID(ctx context.Context, id string) (*domain.ShortURL, error) {
	// filter our ids to find ours
	filter := bson.D{{Key: "_id", Value: id}}
	var shortUrl domain.ShortURL

	// set id to our json struct
	err := r.db.FindOne(ctx, filter).Decode(&shortUrl)
	if err != nil {
		return &shortUrl, err
	}
	fmt.Println(shortUrl)
	return &shortUrl, nil
}

func (r *repository) Update(ctx context.Context, shortUrl *domain.ShortURL) error {
	filter := bson.D{{Key: "_id", Value: shortUrl.Id}}
	updatedRequest, err := r.db.ReplaceOne(ctx, filter, shortUrl)
	if err != nil {
		return err
	}

	if updatedRequest.MatchedCount == 0 {
		return pkg.ErrNotFound
	}
	return nil
}
