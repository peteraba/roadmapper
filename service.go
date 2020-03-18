package main

import (
	"context"
	"encoding/xml"
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
	defaultSvgWidth        = 800
	defaultSvgHeaderHeight = 80
	defaultSvgLineHeight   = 40
)

func Serve(port uint, certFile, keyFile string, rw DbReadWriter, cb CodeBuilder, matomoDomain string, selfHosted bool) {
	// Setup
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{}))
	e.Logger.SetLevel(log.INFO)

	e.File("/favicon.ico", "static/favicon.ico")
	e.Static("/static", "static")
	e.Static("/static", "static")

	e.GET("/", createGetRoadmap(rw, cb, matomoDomain, selfHosted))
	e.POST("/", createPostRoadmap(rw, cb))
	e.GET("/:identifier/svg", createGetRoadRoadmapSVG(rw, cb))
	e.GET("/:identifier", createGetRoadmap(rw, cb, matomoDomain, selfHosted))
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

func createGetRoadRoadmapSVG(rw DbReadWriter, cb CodeBuilder) func(c echo.Context) error {
	return func(c echo.Context) error {
		fw, err := strconv.ParseUint(c.QueryParam("width"), 10, 64)
		if err != nil {
			fw = defaultSvgWidth
		}

		hh, err := strconv.ParseUint(c.QueryParam("height"), 10, 64)
		if err != nil {
			hh = defaultSvgHeaderHeight
		}

		lh, err := strconv.ParseUint(c.QueryParam("lineHeight"), 10, 64)
		if err != nil {
			lh = defaultSvgLineHeight
		}

		fw, hh, lh = getSvgSizes(fw, hh, lh)

		lines, dateFormat, baseUrl, code, err := load(rw, cb, c.Param("identifier"))
		if err != nil {
			log.Print(err)
			return c.HTML(code, fmt.Sprintf("%v", err))
		}

		roadmap, err := linesToRoadmap(lines, dateFormat, baseUrl)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		}

		svg := createSvg(roadmap, float64(fw), float64(hh), float64(lh), dateFormat)

		c.Response().Header().Set(echo.HeaderContentType, "image/svg+xml")

		return c.XML(http.StatusOK, svg)
	}
}

func createGetRoadmap(rw DbReadWriter, cb CodeBuilder, matomoDomain string, selfHosted bool) func(c echo.Context) error {
	return func(c echo.Context) error {
		lines, dateFormat, baseUrl, code, err := load(rw, cb, c.Param("identifier"))
		if err != nil {
			return c.HTML(code, fmt.Sprintf("%v", err))
		}

		roadmap, err := linesToRoadmap(lines, dateFormat, baseUrl)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		}

		output, err := bootstrapRoadmap(roadmap, lines, matomoDomain, dateFormat, baseUrl, selfHosted)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusMethodNotAllowed, fmt.Sprintf("%v", err))
		}

		return c.HTML(http.StatusOK, output)
	}
}

func createPostRoadmap(rw DbReadWriter, cb CodeBuilder) func(c echo.Context) error {
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
		dateFormat := c.FormValue("dateFormat")
		baseUrl := c.FormValue("baseUrl")

		newCode, err := rw.Write(cb, code, content, dateFormat, baseUrl)
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

func Render(rw FileReadWriter, input, output, dateFormat, baseUrl string, fw, hh, lh uint64) error {
	lines, err := rw.Read(input)
	if err != nil {
		return err
	}

	roadmap, err := linesToRoadmap(lines, dateFormat, baseUrl)
	if err != nil {
		return err
	}

	fw, hh, lh = getSvgSizes(fw, hh, lh)

	svg := createSvg(roadmap, float64(fw), float64(hh), float64(lh), dateFormat)

	b, err := xml.Marshal(svg)
	if err != nil {
		return err
	}

	err = rw.Write(output, string(b))

	return err
}

func getSvgSizes(fw, hh, lh uint64) (uint64, uint64, uint64) {
	if fw < defaultSvgWidth {
		fw = defaultSvgWidth
	}
	if hh < defaultSvgHeaderHeight {
		hh = defaultSvgHeaderHeight
	}
	if lh < defaultSvgLineHeight {
		lh = defaultSvgLineHeight
	}

	return fw, hh, lh
}

func load(rw DbReadWriter, cb CodeBuilder, identifier string) ([]string, string, string, int, error) {
	if identifier == "" {
		return []string{}, "", "", 0, nil
	}

	code, err := cb.NewFromString(identifier)
	if err != nil {
		return nil, "", "", http.StatusBadRequest, err
	}

	lines, dateFormat, baseUrl, err := rw.Read(code)
	if err != nil {
		return nil, "", "", http.StatusNotFound, err
	}

	return lines, dateFormat, baseUrl, 0, nil
}

func linesToRoadmap(lines []string, dateFormat, baseUrl string) (Project, error) {
	roadmap, err := parseRoadmap(lines, dateFormat, baseUrl)
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
