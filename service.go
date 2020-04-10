package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/tdewolff/canvas"

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

type renderType string

const (
	svgType renderType = "svg"
)

func Serve(quit chan os.Signal, port uint, certFile, keyFile string, rw DbReadWriter, cb CodeBuilder, matomoDomain, docBaseUrl string, selfHosted bool) {
	// Setup
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{}))
	e.Logger.SetLevel(log.INFO)

	e.File("/favicon.ico", "static/favicon.ico")
	e.Static("/static", "static")
	e.Static("/static", "static")

	e.GET("/", createGetRoadmap(rw, cb, matomoDomain, docBaseUrl, selfHosted))
	e.POST("/", createPostRoadmap(rw, cb))
	e.GET("/:identifier/svg", createGetRoadRoadmapSVG(rw, cb))
	e.GET("/:identifier", createGetRoadmap(rw, cb, matomoDomain, docBaseUrl, selfHosted))
	e.POST("/:identifier", createPostRoadmap(rw, cb))

	// Start server
	go func() {
		if err := startServer(e, port, certFile, keyFile); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func createGetRoadRoadmapSVG(rw DbReadWriter, cb CodeBuilder) func(c echo.Context) error {
	return func(ctx echo.Context) error {
		fw, err := strconv.ParseUint(ctx.QueryParam("width"), 10, 64)
		if err != nil {
			fw = defaultSvgWidth
		}

		hh, err := strconv.ParseUint(ctx.QueryParam("height"), 10, 64)
		if err != nil {
			hh = defaultSvgHeaderHeight
		}

		lh, err := strconv.ParseUint(ctx.QueryParam("lineHeight"), 10, 64)
		if err != nil {
			lh = defaultSvgLineHeight
		}

		fw, hh, lh = getCanvasSizes(fw, hh, lh)

		roadmap, code, err := load(rw, cb, ctx.Param("identifier"))
		if err != nil {
			log.Print(err)
			return ctx.HTML(code, fmt.Sprintf("%v", err))
		}

		cvs := roadmap.ToVisual().Draw(float64(fw), float64(hh), float64(lh))

		img := renderImg(cvs, svgType)

		ctx.Response().Header().Set(echo.HeaderContentType, "image/svg+xml")

		return ctx.String(http.StatusOK, string(img))
	}
}

func createGetRoadmap(rw DbReadWriter, cb CodeBuilder, matomoDomain, docBaseUrl string, selfHosted bool) func(c echo.Context) error {
	return func(ctx echo.Context) error {
		roadmap, code, err := load(rw, cb, ctx.Param("identifier"))
		if err != nil {
			return ctx.HTML(code, err.Error())
		}

		output, err := bootstrapRoadmap(roadmap, matomoDomain, docBaseUrl, selfHosted)
		if err != nil {
			log.Print(err)
			return ctx.HTML(ErrorToHttpCode(err, http.StatusInternalServerError), err.Error())
		}

		return ctx.HTML(http.StatusOK, output)
	}
}

func createPostRoadmap(rw DbReadWriter, cb CodeBuilder) func(c echo.Context) error {
	return func(ctx echo.Context) error {
		prevID, err := getPrevID(cb, ctx.Param("identifier"))
		if err != nil {
			log.Print(err)
			return ctx.Redirect(http.StatusSeeOther, "/?error="+url.QueryEscape(err.Error()))
		}

		content := ctx.FormValue("txt")
		dateFormat := ctx.FormValue("dateFormat")
		baseUrl := ctx.FormValue("baseUrl")
		now := time.Now()

		roadmap := Content(content).ToRoadmap(newCode64().ID(), prevID, dateFormat, baseUrl, now)

		err = rw.Write(cb, roadmap)
		if err != nil {
			log.Print(err)
			return ctx.HTML(http.StatusMethodNotAllowed, err.Error())
		}

		code, err := cb.NewFromID(roadmap.ID)
		if err != nil {
			log.Print(err)
			return ctx.HTML(http.StatusMethodNotAllowed, err.Error())
		}

		newURL := fmt.Sprintf("/%s#", code.String())

		return ctx.Redirect(http.StatusSeeOther, newURL)
	}
}

func getPrevID(cb CodeBuilder, identifier string) (*uint64, error) {
	if identifier == "" {
		return nil, nil
	}

	code, err := cb.NewFromString(identifier)
	if err != nil {
		return nil, HttpError{error: err, status: http.StatusNotFound}
	}

	n := code.ID()

	return &n, err
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

func Render(rw FileReadWriter, content, output, dateFormat, baseUrl string, fw, hh, lh uint64) error {
	fw, hh, lh = getCanvasSizes(fw, hh, lh)

	roadmap := Content(content).ToRoadmap(0, nil, dateFormat, baseUrl, time.Now())
	cvs := roadmap.ToVisual().Draw(float64(fw), float64(hh), float64(lh))
	img := renderImg(cvs, svgType)

	err := rw.Write(output, string(img))

	return err
}

func renderImg(cvs *canvas.Canvas, renderType renderType) []byte {
	var buf bytes.Buffer

	switch renderType {
	case svgType:
		svg := canvas.NewSVG(&buf, cvs.W, cvs.H)

		cvs.Render(svg)
		svg.Close()
	}

	return buf.Bytes()
}

func getCanvasSizes(fw, hh, lh uint64) (uint64, uint64, uint64) {
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

func load(rw DbReadWriter, cb CodeBuilder, identifier string) (*Roadmap, int, error) {
	if identifier == "" {
		return nil, 0, nil
	}

	code, err := cb.NewFromString(identifier)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	roadmap, err := rw.Read(code)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return roadmap, 0, nil
}

func Random(cb CodeBuilder, count int) error {
	var codes []Code
	for i := 0; i < count; i++ {
		codes = append(codes, cb.New())
	}

	displayCodes(codes...)

	return nil
}

func Convert(cb CodeBuilder, id uint64, code string) error {
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
