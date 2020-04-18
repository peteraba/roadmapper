//go:generate go-bindata -o pkg/bindata.go res/migrations/ res/templates/ res/fonts/...

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/peteraba/roadmapper/cmd/roadmapper"
	"github.com/peteraba/roadmapper/pkg/code"
	"github.com/peteraba/roadmapper/pkg/migrations"
	"github.com/peteraba/roadmapper/pkg/roadmap"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var appName = "roadmapper"
var appVersion = "development"
var tag = ""

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	logger, _ := zap.NewProduction()
	defer logger.Sync() // nolint

	cb := code.NewCodeBuilder()

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
					&cli.StringFlag{Name: "matomoDomain", Usage: "matomo domain", EnvVars: []string{"MATOMO_DOMAIN"}},
					&cli.StringFlag{Name: "docBaseUrl", Usage: "documentation base URL", EnvVars: []string{"DOC_BASE_URL"}, Value: "https://docs.rdmp.app"},
					&cli.BoolFlag{Name: "selfHosted", Usage: "self hosted", EnvVars: []string{"SELF_HOSTED"}, Value: false},
					&cli.BoolFlag{Name: "logDbQueries", Usage: "log DB queries", EnvVars: []string{"LOG_DB_QUERIES"}, Value: false},
				},
				Action: func(c *cli.Context) error {
					quit := make(chan os.Signal, 1)
					rw := roadmap.CreateDbReadWriter(
						appName,
						c.String("dbHost"),
						c.String("dbPort"),
						c.String("dbName"),
						c.String("dbUser"),
						c.String("dbPass"),
						c.Bool("logDbQueries"),
					)
					h := roadmap.NewHandler(logger, rw, cb, appVersion, c.String("matomoDomain"), c.String("docBaseUrl"), c.Bool("selfHosted"))
					roadmapper.Serve(quit, c.Uint("port"), c.String("cert"), c.String("key"), h)

					return nil
				},
			},
			{
				Name:    "cli",
				Aliases: []string{"c"},
				Usage:   "renders a roadmap",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "input", Usage: "input file", Aliases: []string{"i"}},
					&cli.StringFlag{Name: "output", Usage: "output file", Aliases: []string{"o"}},
					&cli.StringFlag{Name: "formatFile", Usage: "image format to be used (supported: svg, png, pdf", Aliases: []string{"f"}, Value: "svg", EnvVars: []string{"IMAGE_FORMAT"}},
					&cli.Uint64Flag{Name: "width", Usage: "width of output file", Aliases: []string{"w"}},
					&cli.Uint64Flag{Name: "lineHeight", Usage: "width of output file", Aliases: []string{"lh"}},
					&cli.StringFlag{Name: "dateFormat", Usage: "date format to use", Value: "2006-01-02", EnvVars: []string{"DATE_FORMAT"}},
					&cli.StringFlag{Name: "baseURL", Usage: "base url to use for non-color, non-date extra values", Value: "", EnvVars: []string{"BASE_URL"}},
				},
				Action: func(c *cli.Context) error {
					rw := roadmap.CreateFileReadWriter()
					err := roadmapper.Render(
						rw,
						logger,
						c.String("input"),
						c.String("output"),
						c.String("formatFile"),
						c.String("dateFormat"),
						c.String("baseURL"),
						c.Uint64("width"),
						c.Uint64("lineHeight"),
					)
					if err != nil {
						logger.Error("failed to render roadmap", zap.Error(err))
					}

					return err
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "display appVersion",
				Action: func(c *cli.Context) error {
					fmt.Println(appName, appVersion, tag)

					return nil
				},
			},
			{
				Name:    "migrate:up",
				Aliases: []string{"mu"},
				Usage:   "run migrations",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "dbHost", Usage: "database host", Value: "localhost", EnvVars: []string{"DB_HOST"}},
					&cli.StringFlag{Name: "dbPort", Usage: "database port", Value: "5432", EnvVars: []string{"DB_PORT"}},
					&cli.StringFlag{Name: "dbName", Usage: "database name", Value: "rdmp", EnvVars: []string{"DB_NAME"}},
					&cli.StringFlag{Name: "dbUser", Usage: "database user", Value: "rdmp", EnvVars: []string{"DB_USER"}},
					&cli.StringFlag{Name: "dbPass", Usage: "database password", Value: "", EnvVars: []string{"DB_PASS"}},
					&cli.UintFlag{Name: "steps", Aliases: []string{"s"}, Usage: "number of steps to migrate up"},
				},
				Action: func(c *cli.Context) error {
					n, err := migrations.Up(
						c.String("dbUser"),
						c.String("dbPass"),
						c.String("dbHost"),
						c.String("dbPort"),
						c.String("dbName"),
						c.Int("steps"),
					)

					if err != nil {
						return fmt.Errorf("migration failed: %w", err)
					}

					logger.Info("up migrations run", zap.Int("count", n))

					return nil
				},
			},
			{
				Name:    "migrate:down",
				Aliases: []string{"md"},
				Usage:   "revert migrations",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "dbHost", Usage: "database host", Value: "localhost", EnvVars: []string{"DB_HOST"}},
					&cli.StringFlag{Name: "dbPort", Usage: "database port", Value: "5432", EnvVars: []string{"DB_PORT"}},
					&cli.StringFlag{Name: "dbName", Usage: "database name", Value: "rdmp", EnvVars: []string{"DB_NAME"}},
					&cli.StringFlag{Name: "dbUser", Usage: "database user", Value: "rdmp", EnvVars: []string{"DB_USER"}},
					&cli.StringFlag{Name: "dbPass", Usage: "database password", Value: "", EnvVars: []string{"DB_PASS"}},
					&cli.UintFlag{Name: "steps", Aliases: []string{"s"}, Usage: "number of steps to migrate down"},
				},
				Action: func(c *cli.Context) error {
					n, err := migrations.Down(
						c.String("dbUser"),
						c.String("dbPass"),
						c.String("dbHost"),
						c.String("dbPort"),
						c.String("dbName"),
						c.Int("steps"),
					)

					if err != nil {
						return fmt.Errorf("migration failed: %w", err)
					}

					logger.Info("down migrations run", zap.Int("count", n))

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal("app run error", zap.Error(err))
	}
}
