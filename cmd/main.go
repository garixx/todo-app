package main

import (
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github/garixx/todo-app"
	"github/garixx/todo-app/pkg/handler"
	"github/garixx/todo-app/pkg/repository"
	"github/garixx/todo-app/pkg/service"
	"os"
	"os/signal"
	"syscall"
)

// @title Todo App API
// @version 1.0
// description API server for todo app

// @host localhost:8087
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("config parsing failed:%s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables:%s", err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("postgres.host"),
		Port:     viper.GetString("postgres.port"),
		Username: viper.GetString("postgres.username"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("postgres.dbname"),
		SSLMode:  viper.GetString("postgres.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("DB connection failed:%s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	go func() {
		if err := srv.Run(viper.GetString("rest.port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while starting http server:%s", err.Error())
		}
	}()
	logrus.Print("App started ...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App stopped ...")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shurdown: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on DB shurdown: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
