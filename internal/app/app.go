package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
	"github.com/riabkovK/microgreens/internal/config"
	"github.com/riabkovK/microgreens/internal/handler"
	"github.com/riabkovK/microgreens/internal/service"
	"github.com/riabkovK/microgreens/internal/storage"
	"github.com/riabkovK/microgreens/pkg/auth"
	"github.com/riabkovK/microgreens/pkg/hash"
	"github.com/sirupsen/logrus"
)

func Run() {
	// Config
	cfg, err := config.InitConfig()
	if err != nil {
		logrus.Fatalf("error initializing configs: %v", err)
	}

	// Dependencies
	db, err := storage.NewPostgresDB(cfg)
	if err != nil {
		logrus.Fatalf("error initializing db: %v", err)
	}

	hasher := hash.NewSHA256Hasher(cfg.Auth.PasswordSalt)
	tokenManager, err := auth.NewJWTManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logrus.Fatal(err)
	}
	storages := storage.NewSQLStorage(db)

	deps := service.NewDeps(storages, hasher, tokenManager, cfg)

	// Services, API Handlers
	services := service.NewService(deps)

	handlers := handler.NewHandler(services, tokenManager)

	// Fiber server
	app := fiber.New(fiber.Config{
		EnablePrintRoutes: cfg.Server.EnablePrintRoutes,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
	})
	app.Use(logger.New())

	handlers.SetupRoutes(app)

	go func() {
		logrus.Fatal(app.Listen(":" + cfg.Server.Port))
	}()

	logrus.Print("Microgreens Web App started")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Microgreens Web App shutting down")

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err = app.ShutdownWithContext(ctx); err != nil {
		logrus.Errorf("error occured on server shutting down: %v", err)
	}

	if err = db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %v", err)
	}
}
