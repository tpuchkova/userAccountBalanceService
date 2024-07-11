package main

import (
	"log"
	"os"
	"time"

	"userAccountBalanceService"
	"userAccountBalanceService/pkg/handler"
	"userAccountBalanceService/pkg/repository"
	"userAccountBalanceService/pkg/service"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	db := initializeDatabase()
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	initializeTicker(services)

	srv := new(userAccountBalanceService.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

}

func initializeTicker(services *service.Service) {
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
}

func initializeDatabase() *sqlx.DB {
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
	return db
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
