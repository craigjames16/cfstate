package main

import (
	"log"
	"os"

	"github.com/craigjames16/cfstate/state"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "state",
				Usage: "Commands to manage state",
				Subcommands: []*cli.Command{
					{
						Name:    "check",
						Aliases: []string{"c"},
						Usage:   "Check state",
						Action:  state.CheckStatus,
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
				Name:  "app",
				Usage: "add an app to the list",
				Subcommands: []*cli.Command{
					{
						Name:    "add",
						Aliases: []string{"a"},
						Usage:   "Add application to state",
						Action:  state.AddApp,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "app-name",
								Usage: "Application name",
							},
							&cli.StringFlag{
								Name:  "template-location",
								Usage: "URL relative to repo base for template file",
							},
							&cli.StringFlag{
								Name:  "config-location",
								Usage: "URL relative to repo base for parameter file",
							},
							&cli.StringFlag{
								Name:  "repo",
								Usage: "The name of the repo where the app is located",
							},
						},
					},
				},
			},
			{
				Name:  "repo",
				Usage: "Repo management",
				Subcommands: []*cli.Command{
					{
						Name:    "add",
						Aliases: []string{"r"},
						Usage:   "Add repo to state",
						Action:  state.AddRepo,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "repo-name",
								Usage: "Application name",
							},
							&cli.StringFlag{
								Name:  "repo-url",
								Usage: "URL to repo",
							},
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
