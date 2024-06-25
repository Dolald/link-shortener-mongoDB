package app

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Run(ctx context.Context) error {
	// connect to our mongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		logrus.Error("an error has occurred")
		return err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			logrus.Error(err)
		}
	}()

	return nil
}
