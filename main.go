package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	cli "github.com/urfave/cli/v2"
)

var name = "roadmapper"
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
					&cli.UintFlag{Name: "port", Aliases: []string{"p"}, Value: 0, EnvVars: []string{"PORT"}},
					&cli.StringFlag{Name: "cert", Aliases: []string{"c"}, EnvVars: []string{"SSH_CERT"}},
					&cli.StringFlag{Name: "key", Aliases: []string{"k"}, EnvVars: []string{"SSH_KEY"}},
					&cli.StringFlag{Name: "dbHost", Usage: "database host", Value: "localhost", EnvVars: []string{"DB_HOST"}},
					&cli.StringFlag{Name: "dbPort", Usage: "database port", Value: "5432", EnvVars: []string{"DB_PORT"}},
					&cli.StringFlag{Name: "dbName", Usage: "database name", Value: "rdmp", EnvVars: []string{"DB_NAME"}},
					&cli.StringFlag{Name: "dbUser", Usage: "database user", Value: "rdmp", EnvVars: []string{"DB_USER"}},
					&cli.StringFlag{Name: "dbPass", Usage: "database password", Value: "", EnvVars: []string{"DB_PASS"}},
					&cli.StringFlag{Name: "dateFormat", Usage: "date format", Value: "2006-01-02", EnvVars: []string{"DATE_FORMAT"}},
				},
				Action: func(c *cli.Context) error {
					rw := CreateReadWriter(
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
					&cli.StringFlag{Name: "identifier", Aliases: []string{"r"}, Usage: "roadmap identifier"},
					&cli.StringFlag{Name: "dbHost", Usage: "database host", Value: "localhost", EnvVars: []string{"DB_HOST"}},
					&cli.StringFlag{Name: "dbPort", Usage: "database port", Value: "5432", EnvVars: []string{"DB_PORT"}},
					&cli.StringFlag{Name: "dbName", Usage: "database name", Value: "rdmp", EnvVars: []string{"DB_NAME"}},
					&cli.StringFlag{Name: "dbUser", Usage: "database user", Value: "rdmp", EnvVars: []string{"DB_USER"}},
					&cli.StringFlag{Name: "dbPass", Usage: "database password", Value: "", EnvVars: []string{"DB_PASS"}},
					&cli.StringFlag{Name: "dateFormat", Usage: "dateFormat", Value: "2006-01-02", EnvVars: []string{"DATE_FORMAT"}},
				},
				Action: func(c *cli.Context) error {
					rw := CreateReadWriter(
						c.String("dbHost"),
						c.String("dbPort"),
						c.String("dbName"),
						c.String("dbUser"),
						c.String("dbPass"),
					)
					output, err := Render(rw, cb, c.String("identifier"), c.String("dateFormat"))
					if err != nil {
						log.Print(err)
						return err
					}

					fmt.Println(output)

					return nil
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
					fmt.Println(name, version, tag)
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
