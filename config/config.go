package config

import (
	"context"
	"log"
	"sync"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	RabbitMQURL  string
	MongoURI     string
	DatabaseName string
	ServerPort   string
}

var (
	cfg    *Config
	once   sync.Once
	client *mongo.Client
)

func LoadConfig() *Config {
	once.Do(func() {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		cfg = &Config{
			RabbitMQURL:  viper.GetString("RABBITMQ_URL"),
			MongoURI:     viper.GetString("MONGO_URI"),
			DatabaseName: viper.GetString("DATABASE_NAME"),
			ServerPort:   viper.GetString("SERVER_PORT"),
		}
	})
	return cfg
}

func GetMongoClient() *mongo.Client {
	if client != nil {
		return client
	}

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	return client
}
