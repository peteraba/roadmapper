package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gosuri/uitable"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func Serve(port uint, certFile, keyFile string, rw ReadWriter, cb CodeBuilder) {
	// Setup
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{}))
	e.Logger.SetLevel(log.INFO)

	e.Static("/static", "static")

	e.GET("/:identifier", createGetRoadmap(rw, cb))
	e.POST("/:identifier", createPostRoadmap(rw, cb))

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

func createGetRoadmap(rw ReadWriter, cb CodeBuilder) func(c echo.Context) error {
	return func(c echo.Context) error {
		identifier := c.Param("identifier")

		code, err := cb.NewFromString(identifier)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusBadRequest, fmt.Sprintf("%v", err))
		}

		lines, err := rw.Read(code)
		if err != nil {
			log.Print(err)
			// TODO: Proper 404 check
			// TODO: Proper 404 message
			return c.HTML(http.StatusNotFound, fmt.Sprintf("%v", err))
		}

		output, err := html(lines)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusMethodNotAllowed, fmt.Sprintf("%v", err))
		}

		return c.HTML(http.StatusOK, output)
	}
}

func createPostRoadmap(rw ReadWriter, cb CodeBuilder) func(c echo.Context) error {
	return func(c echo.Context) error {
		identifier := c.Param("identifier")

		code, err := cb.NewFromString(identifier)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusBadRequest, fmt.Sprintf("%v", err))
		}

		content := c.FormValue("roadmap")

		err = rw.Write(cb, code, content)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusMethodNotAllowed, fmt.Sprintf("%v", err))
		}

		return c.Redirect(http.StatusSeeOther, c.Request().URL.String())
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

func Render(rw ReadWriter, cb CodeBuilder, identifier string) error {
	code, err := cb.NewFromString(identifier)
	if err != nil {
		log.Print(err)
		return err
	}

	lines, err := rw.Read(code)
	if err != nil {
		log.Print(err)
		return err
	}

	output, err := html(lines)
	if err != nil {
		log.Print(err)
		return err
	}

	fmt.Println(output)

	return nil
}

func html(lines []string) (string, error) {
	roadmap, err := parseRoadmap(lines)
	if err != nil {
		return "", err
	}

	r := roadmap.ToPublic(roadmap.GetFrom(), roadmap.GetTo())

	return bootstrapRoadmap(r, lines)
}

func Random(cb CodeBuilder, count int) error {
	var codes []Code
	for i := 0; i < count; i++ {
		codes = append(codes, cb.New())
	}

	displayCodes(codes...)

	return nil
}

func Convert(cb CodeBuilder, id int64, code string) error {
	var n Code
	var err error

	if id > 0 {
		n, err = cb.NewFromID(id)
	} else if code != "" {
		n, err = cb.NewFromString(code)
	}

	if err != nil {
		log.Print(err)
		return err
	}

	displayCodes(n)

	return nil
}

func displayCodes(codes ...Code) {
	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true // wrap columns

	table.AddRow("ID", "CODE")
	for _, code := range codes {
		table.AddRow(code.ID(), code.String())
	}

	fmt.Println(table)
}
