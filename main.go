package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	cli "github.com/urfave/cli/v2"
)

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
					&cli.UintFlag{Name: "port", Aliases: []string{"p"}},
					&cli.StringFlag{Name: "cert", Aliases: []string{"c"}},
					&cli.StringFlag{Name: "key", Aliases: []string{"k"}},
					&cli.StringFlag{Name: "dbHost", Usage: "database host", Value: "localhost", EnvVars: []string{"DB_HOST"}},
					&cli.StringFlag{Name: "dbPort", Usage: "database port", Value: "5432", EnvVars: []string{"DB_PORT"}},
					&cli.StringFlag{Name: "dbName", Usage: "database name", Value: "rdmp", EnvVars: []string{"DB_NAME"}},
					&cli.StringFlag{Name: "dbUser", Usage: "database user", Value: "rdmp", EnvVars: []string{"DB_USER"}},
					&cli.StringFlag{Name: "dbPass", Usage: "database password", Value: "", EnvVars: []string{"DB_PASS"}},
				},
				Action: func(c *cli.Context) error {
					rw := CreateReadWriter(c.String("dbHost"), c.String("dbPort"), c.String("dbName"), c.String("dbUser"), c.String("dbPass"))
					Serve(0, "", "", rw, cb)
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
				},
				Action: func(c *cli.Context) error {
					rw := CreateReadWriter(c.String("dbHost"), c.String("dbPort"), c.String("dbName"), c.String("dbUser"), c.String("dbPass"))
					err := Render(rw, cb, c.String("identifier"))
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
