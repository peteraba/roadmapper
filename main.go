package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	cli "gopkg.in/urfave/cli.v2"
)

const inputFile = "roadmap.txt"

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "start server",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "roadmap", Aliases: []string{"r"}, Usage: "path to the roadmap file", DefaultText: "roadmap.txt"},
					&cli.UintFlag{Name: "port", Aliases: []string{"p"}},
					&cli.StringFlag{Name: "cert", Aliases: []string{"c"}},
					&cli.StringFlag{Name: "key", Aliases: []string{"k"}},
				},
				Action: func(c *cli.Context) error {
					serve(0, "", "")
					return nil
				},
			},
			{
				Name:    "cli",
				Aliases: []string{"c"},
				Usage:   "render static assets",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "roadmap", Aliases: []string{"r"}, Usage: "path to the roadmap file", DefaultText: "roadmap.txt"},
				},
				Action: func(c *cli.Context) error {
					roadmap, err := createRoadmap(inputFile)
					if err != nil {
						return err
					}

					output, err := bootstrapRoadmap(roadmap)
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
