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

func Serve(port uint, certFile, keyFile string, rw ReadWriter) {
	// Setup
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{}))
	e.Logger.SetLevel(log.INFO)

	e.Static("/static", "static")

	e.GET("/:identifier", func(c echo.Context) error {
		code, err := NewCode64FromString(c.Param("identifier"))
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
	})

	var putRoadmap = func(c echo.Context, identifier string) error {
		code, err := NewCode64FromString(identifier)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusBadRequest, fmt.Sprintf("%v", err))
		}

		content := c.FormValue("roadmap")

		err = rw.Write(code, content)
		if err != nil {
			log.Print(err)
			return c.HTML(http.StatusMethodNotAllowed, fmt.Sprintf("%v", err))
		}

		return c.Redirect(http.StatusSeeOther, c.Request().URL.String())
	}

	e.POST("/:identifier", func(c echo.Context) error {
		m := c.FormValue("_method")

		if m == "PUT" {
			return putRoadmap(c, c.Param("identifier"))
		}

		return c.HTML(http.StatusMethodNotAllowed, "")
	})

	e.PUT("/:identifier", func(c echo.Context) error {
		return putRoadmap(c, c.Param("identifier"))
	})

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

func Render(rw ReadWriter, identifier string) error {
	code, err := NewCode64FromString(identifier)
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

func Random(count int) error {
	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true // wrap columns

	table.AddRow("ID", "CODE")
	for i := 0; i < count; i++ {
		n := NewCode64()

		table.AddRow(int64(n), n.String())
	}
	fmt.Println(table)

	return nil
}

func Convert(id int, code string) error {
	var n Code64
	var err error

	if id > 0 {
		n = Code64(id)
	} else if code != "" {
		n, err = NewCode64FromString(code)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true // wrap columns

	table.AddRow("ID", "CODE")
	table.AddRow(int64(n), n.String())
	fmt.Println(table)

	return nil
}
