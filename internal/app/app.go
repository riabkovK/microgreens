package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/riabkovK/microgreens/internal/config"
)

type Server struct {
	fiberServer *fiber.App
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		fiberServer: fiber.New(fiber.Config{
			EnablePrintRoutes: cfg.Server.EnablePrintRoutes,
			ReadBufferSize:    cfg.Server.MaxHeaderMegabytes,
			ReadTimeout:       cfg.Server.ReadTimeout,
			WriteTimeout:      cfg.Server.WriteTimeout,
		}),
	}
}

func (s *Server) Run(cfg *config.Config) error {
	return s.fiberServer.Listen(":" + cfg.Server.Port)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.fiberServer.ShutdownWithContext(ctx)
}
