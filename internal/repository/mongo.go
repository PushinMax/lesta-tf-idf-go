package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"os"
)

type MongoConfig struct {
	Port string
	Host string
}

func NewMongoDB(cfg MongoConfig) (*mongo.Client, error) {
	username := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	password := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	host := os.Getenv("MONGO_HOST")     
	port := os.Getenv("MONGO_PORT")    
	// database := os.Getenv("MONGO_DB") 
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/my-db?authSource=admin", username, password, host, port)
	client, _ := mongo.Connect(options.Client().ApplyURI(mongoURI).SetAuth(
			options.Credential{
				AuthSource: "admin", 
				Username:   username,
				Password:   password,
			},
		))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}

	indexes := []mongo.IndexModel{
        {
            Keys:    bson.M{"file_id": 1},
            Options: options.Index().SetUnique(true),
        },
        {
            Keys: bson.M{"user_id": 1},
        },
    }
    
    _, err = client.Database("test").Collection("files").Indexes().CreateMany(context.TODO(), indexes)
	if err != nil {
		log.Printf("index didn't create: %s", err.Error())
	}
	return client, nil
}