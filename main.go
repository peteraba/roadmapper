package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "start server",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "roadmap", Aliases: []string{"r"}, Usage: "path to the roadmap file", Value: "roadmap.txt"},
					&cli.UintFlag{Name: "port", Aliases: []string{"p"}},
					&cli.StringFlag{Name: "cert", Aliases: []string{"c"}},
					&cli.StringFlag{Name: "key", Aliases: []string{"k"}},
				},
				Action: func(c *cli.Context) error {
					serve(0, "", "", c.String("roadmap"))
					return nil
				},
			},
			{
				Name:    "cli",
				Aliases: []string{"c"},
				Usage:   "render static assets",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "roadmap", Aliases: []string{"r"}, Usage: "path to the roadmap file", Value: "roadmap.txt"},
				},
				Action: func(c *cli.Context) error {
					output, err := html(c.String("roadmap"))
					if err != nil {
						return err
					}

					fmt.Println(output)

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
