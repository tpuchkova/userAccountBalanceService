package main

import (
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
	"userAccountBalanceService"
	"userAccountBalanceService/pkg/handler"
	"userAccountBalanceService/pkg/repository"
	"userAccountBalanceService/pkg/service"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		Username: os.Getenv("DATABASE_USERNAME"),
		DBName:   os.Getenv("DATABASE_NAME"),
		SSLMode:  "disable",
		Password: os.Getenv("DATABASE_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(userAccountBalanceService.Server)

	ticker := time.NewTicker(time.Duration(viper.GetInt("cancelTransactionIntervalMinutes")) * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				err := services.CancelLatestOddTransactions()
				if err != nil {
					return
				}
			}
		}
	}()

	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
