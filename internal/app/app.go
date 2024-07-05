package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"shortener/configs"
	"shortener/internal/handler_http"
	"shortener/internal/repository"
	"shortener/internal/service"

	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Run(ctx context.Context) {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// load configs variable
	if err := configsInit(); err != nil {
		panic(err)
	}

	// load env variable
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	clientOptions := options.Client().ApplyURI("mongodb://" + viper.GetString("db.host") + "27017").SetAuth(options.Credential{
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	// if request handle  more than 10 seconds cancle the request
	ctx, cancel := context.WithTimeout(context.Background(), configs.ContextWaiting*time.Second)
	defer cancel()

	// connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	// defer disconnecting
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	repository, err := repository.NewRepository(ctx, client)
	if err != nil {
		panic(err)
	}

	service := service.NewService(repository)
	handler := handler_http.NewHandler(service)

	log.Println("The database was created and indices were set up successfully")

	if err = http.ListenAndServe(viper.GetString("db.host")+viper.GetString("port"), handler.InitRoutes()); err != nil {
		panic(err)
	}
}

func configsInit() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
