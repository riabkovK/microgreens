package main

import (
	"github.com/riabkovK/microgreens/internal/storage"
	"github.com/riabkovK/microgreens/pkg/handler"
	"github.com/riabkovK/microgreens/pkg/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %v", err)
	}

	db, err := storage.NewPostgresDB(storage.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("error initializing db: %v", err)
	}

	storages := storage.NewSQLStorage(db)
	services := service.NewService(storages)
	handlers := handler.NewHandler(services)

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: viper.GetBool("app.enablePrintRoutes"),
		ReadBufferSize:    viper.GetInt("app.maxHeaderSize"),
		ReadTimeout:       viper.GetDuration("readTimeout"),
		WriteTimeout:      viper.GetDuration("writeTimeout"),
	})
	app.Use(logger.New(), recover.New())

	handlers.SetupRoutes(app)

	logrus.Fatal(app.Listen(":" + viper.GetString("app.port")))
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Error("read config file err: %v", err)
		return err
	}

	// Add .env (secrets)
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	if err := viper.MergeInConfig(); err != nil {
		logrus.Error("merge .env file err: %v", err)
		return err
	}

	return nil
}
