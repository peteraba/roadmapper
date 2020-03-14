package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	cli "github.com/urfave/cli/v2"
)

var applicationName = "roadmapper"
var tag = ""
var version = "development"

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	cb := NewCodeBuilder()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "start server",
				Flags: []cli.Flag{
					&cli.UintFlag{Name: "port", Usage: "port to be used by the server", Aliases: []string{"p"}, Value: 0, EnvVars: []string{"PORT"}},
					&cli.StringFlag{Name: "cert", Usage: "SSH cert used for https", Aliases: []string{"c"}, EnvVars: []string{"SSH_CERT"}},
					&cli.StringFlag{Name: "key", Usage: "SSH key used for https", Aliases: []string{"k"}, EnvVars: []string{"SSH_KEY"}},
					&cli.StringFlag{Name: "dbHost", Usage: "database host", Value: "localhost", EnvVars: []string{"DB_HOST"}},
					&cli.StringFlag{Name: "dbPort", Usage: "database port", Value: "5432", EnvVars: []string{"DB_PORT"}},
					&cli.StringFlag{Name: "dbName", Usage: "database name", Value: "rdmp", EnvVars: []string{"DB_NAME"}},
					&cli.StringFlag{Name: "dbUser", Usage: "database user", Value: "rdmp", EnvVars: []string{"DB_USER"}},
					&cli.StringFlag{Name: "dbPass", Usage: "database password", Value: "", EnvVars: []string{"DB_PASS"}},
					&cli.StringFlag{Name: "dateFormat", Usage: "date format to use", Value: "2006-01-02", EnvVars: []string{"DATE_FORMAT"}},
				},
				Action: func(c *cli.Context) error {
					rw := CreateDbReadWriter(
						applicationName,
						c.String("dbHost"),
						c.String("dbPort"),
						c.String("dbName"),
						c.String("dbUser"),
						c.String("dbPass"),
					)
					Serve(c.Uint("port"), c.String("cert"), c.String("key"), rw, cb, c.String("dateFormat"))

					return nil
				},
			},
			{
				Name:    "cli",
				Aliases: []string{"c"},
				Usage:   "render static assets",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "input", Usage: "input file", Aliases: []string{"i"}},
					&cli.StringFlag{Name: "output", Usage: "output file", Aliases: []string{"o"}},
					&cli.Uint64Flag{Name: "width", Usage: "width of output file", Aliases: []string{"w"}},
					&cli.Uint64Flag{Name: "headerHeight", Usage: "width of output file", Aliases: []string{"hh"}},
					&cli.Uint64Flag{Name: "lineHeight", Usage: "width of output file", Aliases: []string{"lh"}},
					&cli.StringFlag{Name: "dateFormat", Usage: "date format to use", Value: "2006-01-02", EnvVars: []string{"DATE_FORMAT"}},
				},
				Action: func(c *cli.Context) error {
					rw := CreateFileReadWriter()
					err := Render(
						rw,
						c.String("input"),
						c.String("output"),
						c.String("dateFormat"),
						c.Uint64("width"),
						c.Uint64("headerHeight"),
						c.Uint64("lineHeight"),
					)
					if err != nil {
						log.Print(err)
					}

					return err
				},
			},
			{
				Name:    "random",
				Aliases: []string{"r"},
				Usage:   "generate random numbers",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "count", Aliases: []string{"c"}, Usage: "count of random numbers to generate", Value: 5},
				},
				Action: func(c *cli.Context) error {
					err := Random(cb, c.Int("count"))
					if err != nil {
						log.Print(err)
					}
					return err
				},
			},
			{
				Name:    "convert",
				Aliases: []string{"co"},
				Usage:   "convert between id and code",
				Flags: []cli.Flag{
					&cli.Int64Flag{Name: "id", Aliases: []string{"i"}, Usage: "id to convert to code"},
					&cli.StringFlag{Name: "code", Aliases: []string{"c"}, Usage: "code to convert to id"},
				},
				Action: func(c *cli.Context) error {
					err := Convert(cb, c.Int64("id"), c.String("code"))
					return err
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "display version",
				Action: func(c *cli.Context) error {
					fmt.Println(applicationName, version, tag)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
