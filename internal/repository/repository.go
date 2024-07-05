package repository

import (
	"context"
	"fmt"
	domain "shortener/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Insert(ctx context.Context, shortURL *domain.ShortURL) error
	FindByID(ctx context.Context, id string) (*domain.ShortURL, error)
	Update(ctx context.Context, shortUrl *domain.ShortURL) error
	Delete(ctx context.Context, shortUrl string) error
}

// create new mongo DB collection
type repository struct {
	db *mongo.Collection
}

func (r *repository) Insert(ctx context.Context, shortURL *domain.ShortURL) error {
	if _, err := r.db.InsertOne(ctx, shortURL); err != nil {
		return fmt.Errorf("Insert/InsertOne: %w", err)
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, shortURL string) (*domain.ShortURL, error) {
	// filter our ids to find ours
	filter := bson.D{{Key: "_id", Value: shortURL}}
	var shortUrl domain.ShortURL

	// set id to our json struct
	err := r.db.FindOne(ctx, filter).Decode(&shortUrl)
	if err != nil {
		return &shortUrl, fmt.Errorf("FindByID/Decode: %w", err)
	}

	return &shortUrl, nil
}

func (r *repository) Update(ctx context.Context, shortUrl *domain.ShortURL) error {
	filter := bson.D{{Key: "_id", Value: shortUrl.Id}}
	updateResult, err := r.db.ReplaceOne(ctx, filter, shortUrl)
	if err != nil {
		return fmt.Errorf("Update/ReplaceOne: %w", err)
	}
	// check finding and changing the document
	if updateResult.MatchedCount == 0 {
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, shortUrl string) error {
	filter := bson.D{{Key: "_id", Value: shortUrl}}

	updateResult, err := r.db.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("Delete/DeleteOne: %w", err)
	}
	// check finding and changing the document
	if updateResult.DeletedCount == 0 {
		return err
	}

	return nil
}
