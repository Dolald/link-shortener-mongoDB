package repository

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRepository(ctx context.Context, client *mongo.Client) (Repository, error) {
	repository := &repository{
		// use DB with "mongoDB" name and "sortUrls" collection
		db: client.Database(viper.GetString("db.dbname")).Collection(viper.GetString("db.collection")),
	}

	if err := repository.createIndices(ctx); err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *repository) createIndices(ctx context.Context) error {
	_, err := r.db.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "expireAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	fmt.Println("The database was created")
	return err
}
