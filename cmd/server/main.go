package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PushinMax/lesta-tf-idf-go/internal/handler"
	"github.com/PushinMax/lesta-tf-idf-go/internal/repository"
	"github.com/PushinMax/lesta-tf-idf-go/internal/server"
	"github.com/PushinMax/lesta-tf-idf-go/internal/service"
	"github.com/PushinMax/lesta-tf-idf-go/internal/session"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)



func main() {
	if err := InitConfig(); err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: os.Getenv("DB_USER"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatal(err)
	}
	mongoDB, err := repository.NewMongoDB(repository.MongoConfig{
		Host: "localhost",
		Port: "27017",
	})
	if err != nil {
		log.Fatal(err)
	}
	
	repository := repository.New(db, mongoDB)
	session := session.New()
	
	service := service.New(session, repository)
	handler := handler.New(service)
	server := new(server.Server)


	go func() {
		if err := server.Run(viper.GetString("server.port"), handler.Init()); err != nil {
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	<-ch
	_ = server.Shutdown(context.Background())
	_ = db.Close()
	_ = mongoDB.Client().Disconnect(context.Background())
}

func InitConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(".") 
    viper.SetConfigType("yml")
	return viper.ReadInConfig()
}




