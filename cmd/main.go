package main

import (
	"github.com/riabkovK/microgreens/internal/storage"
	"github.com/riabkovK/microgreens/pkg/handler"
	"github.com/riabkovK/microgreens/pkg/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// TODO Добавить:
// 1) рефреш токен
// 2) сборку бинаря + проброс конфига (соответственно, флаги), поиск переменных из окружения (viper)
// 3) docker compose,
// 4) Тесты
// 5) swagger
// 6) CI?
// 7) UI (на реакте)

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

	go func() {
		logrus.Fatal(app.Listen(":" + viper.GetString("app.port")))
	}()

	logrus.Print("Microgreens Web App started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Microgreens Web App shutting down")
	if err = app.Shutdown(); err != nil {
		logrus.Errorf("error occured on server shutting down: %v", err)
	}

	if err = db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %v", err)
	}
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
