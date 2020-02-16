package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func server(port uint, certFile, keyFile, inputFile string) {
	// Setup
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{}))
	e.Logger.SetLevel(log.INFO)

	e.Static("/static", "static")

	e.GET("/roadmap", func(c echo.Context) error {
		output, err := html(inputFile)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusMethodNotAllowed, fmt.Sprintf("%v", err))
		}

		return c.HTML(http.StatusOK, output)
	})

	var putRoadmap = func(c echo.Context) error {
		content := c.FormValue("roadmap")

		err := save(inputFile, content)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusMethodNotAllowed, fmt.Sprintf("%v", err))
		}

		return c.Redirect(http.StatusSeeOther, "/roadmap")
	}

	e.POST("/roadmap", func(c echo.Context) error {
		m := c.FormValue("_method")

		if m == "PUT" {
			return putRoadmap(c)
		}

		return c.HTML(http.StatusMethodNotAllowed, "")
	})

	e.PUT("/roadmap", putRoadmap)

	// Start server
	go func() {
		if err := startServer(e, port, certFile, keyFile); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func startServer(e *echo.Echo, port uint, certFile, keyFile string) error {
	var minPort, maxPort uint = 1323, 11000

	if port > 0 {
		minPort = port
		maxPort = port
	}

	f := startWrapper(e, certFile, keyFile)

	for p := minPort; p <= maxPort; p++ {
		e.Logger.Info(fmt.Sprintf("trying port: %d", p))
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

func commandLine(inputFile string) error {
	output, err := html(inputFile)
	if err != nil {
		return err
	}

	fmt.Println(output)

	return nil
}

func html(inputFile string) (string, error) {
	lines, err := readRoadmap(inputFile)
	if err != nil {
		return "", err
	}

	roadmap, err := parseRoadmap(lines)
	if err != nil {
		return "", err
	}

	r := roadmap.ToPublic(roadmap.GetFrom(), roadmap.GetTo())

	return bootstrapRoadmap(r, lines)
}

func save(inputFile, content string) error {
	if content == "" {
		return errors.New("content must not be empty.")
	}

	lines := strings.Split(content, "\r\n")

	_, err := parseRoadmap(lines)
	if err != nil {
		return err
	}

	return writeRoadmap(inputFile, content)
}
