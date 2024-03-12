package main

import (
	"context"
	"errors"
	"github.com/Warh40k/cloud-manager/internal/api/cache"
	"github.com/Warh40k/cloud-manager/internal/api/cache/redis_cache"
	httpServ "github.com/Warh40k/cloud-manager/internal/api/handler/http"
	"github.com/Warh40k/cloud-manager/internal/api/repository"
	"github.com/Warh40k/cloud-manager/internal/api/repository/postgres"
	"github.com/Warh40k/cloud-manager/internal/api/service"
	"github.com/Warh40k/cloud-manager/internal/app"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("local")
	return viper.ReadInConfig()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("Ошибка чтения конфигурации: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Ошибка чтения переменных окружения: %s", err.Error())
	}

	log := setupLogger(viper.GetString("env"))

	pgCfg := postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}
	db, err := postgres.NewPostgresDB(pgCfg)
	if err != nil {
		log.Error("Ошибка подключения к базе данных: %s", err.Error())
		return
	}

	//redis
	rdCfg := redis_cache.Config{
		Addr:     viper.GetString("cache.addr"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       viper.GetInt("cache.db"),
	}
	ctx := context.Background()
	rd, err := redis_cache.NewRedisConn(ctx, rdCfg)
	if err != nil {
		log.Error("Ошибка подключения к кэшу: %s", err.Error())
		return
	}

	repos := repository.NewRepository(db, log)
	cacheDb := cache.NewCache(ctx, rd)
	services := service.NewService(repos, cacheDb, log)
	handlers := httpServ.NewHandler(services, log)
	serv := new(app.App)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		if err = serv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Ошибка запуска http сервера: %s", err.Error())
		}
	}()
	log.Info("server started")
	<-quit

	log.Info("trying to gracefull shutdown")
	if err = serv.Shutdown(context.Background()); err != nil {
		log.With(slog.String("err", err.Error())).Error("error occured on server shutting down:")
	}

	if err = db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}

	log.Info("gracefully stopped")
}
