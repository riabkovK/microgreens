package server

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/riabkovK/microgreens/internal/config"
)

type Server struct {
	fiberServer *fiber.App
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		fiberServer: fiber.New(fiber.Config{
			EnablePrintRoutes: cfg.Server.EnablePrintRoutes,
			ReadBufferSize:    cfg.Server.MaxHeaderMegabytes << 10,
			ReadTimeout:       cfg.Server.ReadTimeout,
			WriteTimeout:      cfg.Server.WriteTimeout,
		}),
	}
}

func (s *Server) Run(cfg *config.Config) error {
	s.initDefaultMiddlewares(cfg)
	return s.fiberServer.Listen(":" + cfg.Server.Port)
}

func (s *Server) initDefaultMiddlewares(cfg *config.Config) fiber.Router {
	return s.fiberServer.Use(
		logger.New(logger.Config{
			Format:        cfg.Logger.Format,
			TimeFormat:    cfg.Logger.TimeFormat,
			TimeZone:      cfg.Logger.TimeZone,
			TimeInterval:  cfg.Logger.TimeInterval,
			Output:        cfg.Logger.Output,
			DisableColors: cfg.Logger.DisableColors,
		}),
		recover.New(),
		limiter.New(limiter.Config{
			KeyGenerator: cfg.Limiter.KeyGenerator,
			Max:          cfg.Limiter.Max,
			Expiration:   cfg.Limiter.Expiration,
		}))
}
func (s *Server) Shutdown(ctx context.Context) error {
	return s.fiberServer.ShutdownWithContext(ctx)
}
