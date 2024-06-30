package repository

import (
	"context"
	"fmt"
	"shortener/pkg"

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
