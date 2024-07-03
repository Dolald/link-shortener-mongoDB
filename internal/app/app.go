package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"shortener/internal/handler_http"
	"shortener/internal/repository"
	"shortener/internal/service"

	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Run(ctx context.Context) error {
	// load configs variable
	if err := configsInit(); err != nil {
		return fmt.Errorf("Run.configsInit: %w", err)
	}

	// load env variable
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("app/Run/godotenv: %w", err)
	}

	clientOptions := options.Client().ApplyURI("mongodb://root:" + os.Getenv("DB_PASSWORD") + "qwerty@localhost:27017").SetAuth(options.Credential{ // чекнуть что это
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	// if request handle  more than 10 seconds cancle the request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("Run/Connect: %w", err)
	}
	// defer disconnecting
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(fmt.Errorf("Run/Disconnect: %w", err))
		}
	}()

	repository, err := repository.NewRepository(ctx, client)
	if err != nil {
		log.Fatal(fmt.Errorf("Run/NewRepository: %w", err))
	}

	service := service.NewService(repository)
	handler := handler_http.NewHandler(service)

	log.Println("The database was created and indices were set up successfully")

	return http.ListenAndServe("localhost:8080", handler.InitRoutes())
}

func configsInit() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
