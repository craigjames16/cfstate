package main

import (
	"log"
	"os"

	"github.com/craigjames16/cfstate/state"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "lang",
				Aliases: []string{"l"},
				Value:   "english",
				Usage:   "Language for the greeting",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "state",
				Subcommands: []*cli.Command{
					{
						Name:    "check",
						Aliases: []string{"c"},
						Usage:   "Check state",
						Action:  state.CheckState,
					},
					{
						Name:    "sync",
						Aliases: []string{"s"},
						Usage:   "sync state",
						Action:  state.SyncState,
					},
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a task to the list",
				Action: func(c *cli.Context) error {
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
