package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PushinMax/lesta-tf-idf-go/internal/handler"
	"github.com/PushinMax/lesta-tf-idf-go/internal/server"
	"github.com/PushinMax/lesta-tf-idf-go/internal/service"
	"github.com/PushinMax/lesta-tf-idf-go/internal/session"
	// "github.com/joho/godotenv"
	// "github.com/spf13/viper"
)



func main() {
	/*if err := InitConfig(); err != nil {
		log.Fatal(err)
	}*/
	session := session.New()
	service := service.New(session)
	handler := handler.New(service)
	server := new(server.Server)


	go func() {
		if err := server.Run("8080", handler.Init()); err != nil {
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	<-ch
	_ = server.Shutdown(context.Background())
	// _ = db.Close()

	
	
}

/*func InitConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(".") 
    viper.SetConfigType("yml")
	return viper.ReadInConfig()
}*/




