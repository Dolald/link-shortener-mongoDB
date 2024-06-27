package repository

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DAO - data acsess object - main DB logic
// create new mongo DB collection
type UrlDAO struct {
	db *mongo.Collection
}

func NewUrlDAO(ctx context.Context, client *mongo.Client) (*UrlDAO, error) {
	dao := &UrlDAO{
		// use DB with "mongoDB" name and "sortUrls" collection
		db: client.Database(viper.GetString("db.dbname")).Collection(viper.GetString("db.collection")),
	}

	if err := dao.createIndices(ctx); err != nil {
		return nil, err
	}

	return dao, nil
}

func (dao *UrlDAO) createIndices(ctx context.Context) error {
	_, err := dao.db.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{"expireAt", 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	fmt.Println("The database was created")
	return err
}
