package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"

	rmid "github.com/peteraba/roadmapper/pkg/middleware"
	"github.com/peteraba/roadmapper/pkg/roadmap"
)

type Server struct {
	echo              *echo.Echo
	handler           *roadmap.Handler
	assetsDir         string
	certFile, keyFile string
}

func NewServer(handler *roadmap.Handler, assetsDir, certFile, keyFile string) *Server {
	// Setup
	e := echo.New()

	e.HideBanner = true

	e.Use(rmid.LoggerWithConfig(&rmid.LoggerConfig{
		Logger: handler.Logger,
	}))
	e.Use(middleware.SecureWithConfig(middleware.DefaultSecureConfig))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableStackAll:   true,
		DisablePrintStack: true,
	}))

	e.File("/favicon.ico", fmt.Sprintf("%s/static/favicon.ico", assetsDir))
	e.Static("/static", assetsDir)

	e.GET("/", handler.GetRoadmapHTML)
	e.POST("/", handler.CreateRoadmapHTML)
	e.GET("/:identifier/:format", handler.GetRoadmapImage)
	e.GET("/:identifier", handler.GetRoadmapHTML)
	e.POST("/:identifier", handler.CreateRoadmapHTML)

	apiGroup := e.Group("/api")
	apiGroup.GET("/:identifier", handler.GetRoadmapJSON)
	apiGroup.POST("/", handler.CreateRoadmapJSON)

	return &Server{
		echo:      e,
		handler:   handler,
		assetsDir: assetsDir,
		certFile:  certFile,
		keyFile:   keyFile,
	}
}

func (s *Server) StartWithTeardown(port uint) func() {
	quit := make(chan os.Signal, 1)

	go s.Start(quit, port)

	teardown := func() {
		quit <- os.Interrupt
	}

	return teardown
}

// Serve sets up a Roadmapper HTTP service using echo
func (s *Server) Start(quit chan os.Signal, port uint) {
	// Start server
	go func() {
		if err := s.serveInBackground(port); err != nil {
			s.handler.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.echo.Shutdown(ctx); err != nil {
		s.handler.Logger.Fatal("shutdown error", zap.Error(err))
	}
}

func (s *Server) serveInBackground(port uint) error {
	var minPort, maxPort uint = 1323, 11000

	if port > 0 {
		minPort = port
		maxPort = port
	}

	f := startWrapper(s.echo, s.certFile, s.keyFile)

	for p := minPort; p <= maxPort; p++ {
		s.handler.Logger.Info(fmt.Sprintf("trying port: %d", p))
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
