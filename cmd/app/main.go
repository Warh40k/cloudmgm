package main

import (
	"github.com/Warh40k/cloud-manager/internal/app"
	"github.com/Warh40k/cloud-manager/internal/handler"
	"github.com/Warh40k/cloud-manager/internal/repository"
	"github.com/Warh40k/cloud-manager/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("Ошибка чтения конфигурации: %s", err.Error())
	}

	/*	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Ошибка чтения переменных окружения: %s", err.Error())
	}*/

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("Ошибка подключения к базе данных: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	serv := new(app.Server)
	if err = serv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Ошибка запуска http сервера: %s", err.Error())
	}
}
