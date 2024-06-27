package app

import (
	"context"
	"log"
	"shortener/pkg/repository"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Run(ctx context.Context) error {

	if err := configsInit(); err != nil {
		return err
	}

	clientOptions := options.Client().ApplyURI("mongodb://root:qwerty@localhost:27017").SetAuth(options.Credential{
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключение к MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Отключение от MongoDB при завершении работы
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Создание экземпляра UrlDAO
	_, err = repository.NewUrlDAO(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("The database was created and indices were set up successfully")

	return err
}

func configsInit() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
