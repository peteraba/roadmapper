package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"go.uber.org/zap"

	"github.com/peteraba/roadmapper/pkg/middleware"
	"github.com/peteraba/roadmapper/pkg/roadmap"
)

// Serve sets up a Roadmapper HTTP service using echo
func Serve(quit chan os.Signal, port uint, certFile, keyFile, assetsDir string, h *roadmap.Handler) {
	// Setup
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Logger: h.Logger,
	}))

	e.File("/favicon.ico", fmt.Sprintf("%s/static/favicon.ico", assetsDir))
	e.Static("/static", assetsDir)

	e.GET("/", h.GetRoadmapHTML)
	e.POST("/", h.CreateRoadmapHTML)
	e.GET("/:identifier/:format", h.GetRoadmapImage)
	e.GET("/:identifier", h.GetRoadmapHTML)
	e.POST("/:identifier", h.CreateRoadmapHTML)

	apiGroup := e.Group("/api")
	apiGroup.GET("/:identifier", h.GetRoadmapJSON)
	apiGroup.POST("/", h.CreateRoadmapJSON)

	// Start server
	go func() {
		if err := startServer(e, h.Logger, port, certFile, keyFile); err != nil {
			h.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		h.Logger.Fatal("shutdown error", zap.Error(err))
	}
}

func startServer(e *echo.Echo, l *zap.Logger, port uint, certFile, keyFile string) error {
	var minPort, maxPort uint = 1323, 11000

	if port > 0 {
		minPort = port
		maxPort = port
	}

	f := startWrapper(e, certFile, keyFile)

	for p := minPort; p <= maxPort; p++ {
		l.Info(fmt.Sprintf("trying port: %d", p))
		if err := f(p); err != nil && p >= maxPort {
			return err
		}
	}

	return nil
}

func startWrapper(e *echo.Echo, certFile, keyFile string) func(port uint) error {
	if certFile == "" || keyFile == "" {
		return func(port uint) error {
			return e.Start(fmt.Sprintf(":%d", port))
		}
	}

	return func(port uint) error {
		return e.StartTLS(fmt.Sprintf(":%d", port), certFile, keyFile)
	}
}
