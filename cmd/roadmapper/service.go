package roadmapper

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/peteraba/roadmapper/pkg/roadmap"
	"go.uber.org/zap"
)

// Serve sets up a Roadmapper HTTP service using echo
func Serve(quit chan os.Signal, port uint, certFile, keyFile string, h *roadmap.Handler) {
	// Setup
	e := echo.New()

	e.File("/favicon.ico", "static/favicon.ico")
	e.Static("/static", "static")
	e.Static("/static", "static")

	e.GET("/", h.GetRoadmapHTML)
	e.POST("/", h.CreateRoadmapHTML)
	e.GET("/:identifier/:format", h.GetRoadmapImage)
	e.GET("/:identifier", h.GetRoadmapHTML)
	e.POST("/:identifier", h.CreateRoadmapHTML)

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

// Render renders a roadmap
func Render(rw roadmap.FileReadWriter, l *zap.Logger, content, output string, fileFormat, dateFormat, baseUrl string, fw, lh uint64) error {
	format, err := roadmap.NewFormatType(fileFormat)
	if err != nil {
		l.Info("format is not supported", zap.Error(err))

		return err
	}

	fw, lh = roadmap.GetCanvasSizes(fw, lh)

	r := roadmap.Content(content).ToRoadmap(0, nil, "", dateFormat, baseUrl, time.Now())
	cvs := r.ToVisual().Draw(float64(fw), float64(lh))
	img := roadmap.RenderImg(cvs, format)

	err = rw.Write(output, string(img))

	return err
}
