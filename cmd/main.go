package main

import (
	"awesomeProject"
	"awesomeProject/pkg/handler"
	"awesomeProject/pkg/repository"
	"awesomeProject/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(awesomeProject.Server)

	//if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
	//	log.Fatalf("error occured while running http server: %s", err.Error())
	//}

	//ticker := time.NewTicker(5 * time.Second)
	//
	//go func() {
	//	for {
	//		select {
	//		case <-ticker.C:
	//			err := services.CancelLatestOddTransactions()
	//			if err != nil {
	//				return
	//			}
	//		}
	//	}
	//}()

	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
