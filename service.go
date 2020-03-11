package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gosuri/uitable"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

const (
	defaultSvgWidth        = 915
	defaultSvgHeaderHeight = 80
	defaultSvgLineHeight   = 40
)

func Serve(port uint, certFile, keyFile string, rw ReadWriter, cb CodeBuilder, dateFormat string) {
	// Setup
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{}))
	e.Logger.SetLevel(log.INFO)

	e.File("/favicon.ico", "static/favicon.ico")
	e.Static("/static", "static")
	e.Static("/static", "static")

	e.GET("/", createGetRoadmap(rw, cb, dateFormat))
	e.POST("/", createPostRoadmap(rw, cb))
	e.GET("/:identifier/svg", createGetRoadRoadmapSVG(rw, cb, dateFormat))
	e.GET("/:identifier", createGetRoadmap(rw, cb, dateFormat))
	e.POST("/:identifier", createPostRoadmap(rw, cb))

	// Start server
	go func() {
		if err := startServer(e, port, certFile, keyFile); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func createGetRoadRoadmapSVG(rw ReadWriter, cb CodeBuilder, dateFormat string) func(c echo.Context) error {
	return func(c echo.Context) error {
		fw, err := strconv.ParseInt(c.QueryParam("width"), 10, 64)
		if err != nil {
			fw = defaultSvgWidth
		}

		fh, err := strconv.ParseInt(c.QueryParam("height"), 10, 64)
		if err != nil {
			fh = defaultSvgHeaderHeight
		}

		lh, err := strconv.ParseInt(c.QueryParam("line-height"), 10, 64)
		if err != nil {
			lh = defaultSvgLineHeight
		}

		lines, code, err := load(rw, cb, c.Param("identifier"))
		if err != nil {
			log.Print(err)
			return c.HTML(code, fmt.Sprintf("%v", err))
		}

		roadmap, err := linesToRoadmap(lines, dateFormat)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		}

		svg := createSvg(roadmap, float64(fw), float64(fh), float64(lh), dateFormat)

		return c.XML(http.StatusOK, svg)
	}
}

func createGetRoadmap(rw ReadWriter, cb CodeBuilder, dateFormat string) func(c echo.Context) error {
	return func(c echo.Context) error {
		lines, code, err := load(rw, cb, c.Param("identifier"))
		if err != nil {
			return c.HTML(code, fmt.Sprintf("%v", err))
		}

		roadmap, err := linesToRoadmap(lines, dateFormat)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		}

		output, err := bootstrapRoadmap(roadmap, lines)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusMethodNotAllowed, fmt.Sprintf("%v", err))
		}

		return c.HTML(http.StatusOK, output)
	}
}

func createPostRoadmap(rw ReadWriter, cb CodeBuilder) func(c echo.Context) error {
	return func(c echo.Context) error {
		var (
			code = cb.New()
			err  error
		)

		identifier := c.Param("identifier")

		if identifier != "" {
			code, err = cb.NewFromString(identifier)
			if err != nil {
				log.Print(err)
				return c.HTML(http.StatusBadRequest, fmt.Sprintf("%v", err))
			}
		}

		content := c.FormValue("txt")

		newCode, err := rw.Write(cb, code, content)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusMethodNotAllowed, fmt.Sprintf("%v", err))
		}

		newURL := fmt.Sprintf("/%s", newCode.String())

		return c.Redirect(http.StatusSeeOther, newURL)
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

func Render(rw ReadWriter, cb CodeBuilder, identifier, dateFormat string) (string, error) {
	code, err := cb.NewFromString(identifier)
	if err != nil {
		return "", err
	}

	lines, err := rw.Read(code)
	if err != nil {
		return "", err
	}

	roadmap, err := linesToRoadmap(lines, dateFormat)
	if err != nil {
		return "", err
	}

	output, err := bootstrapRoadmap(roadmap, lines)
	if err != nil {
		return "", err
	}

	return output, nil
}

func load(rw ReadWriter, cb CodeBuilder, identifier string) ([]string, int, error) {
	if identifier == "" {
		return []string{}, 0, nil
	}

	code, err := cb.NewFromString(identifier)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	lines, err := rw.Read(code)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return lines, 0, nil
}

func linesToRoadmap(lines []string, dateFormat string) (Project, error) {
	roadmap, err := parseRoadmap(lines, dateFormat)
	if err != nil {
		return Project{}, err
	}

	return roadmap.ToPublic(roadmap.GetStart(), roadmap.GetEnd()), nil
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
