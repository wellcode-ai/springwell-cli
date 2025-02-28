package main

import (
	"os"

	"github.com/springwell/cli/pkg/commands"
	"github.com/springwell/cli/pkg/util"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "springwell",
		Usage:   "Developer-friendly CLI for Spring Boot",
		Version: "0.1.0",
		Authors: []*cli.Author{
			{
				Name:  "SpringWell Team",
				Email: "info@springwell.dev",
			},
		},
		Commands: []*cli.Command{
			commands.NewProjectCommand(),
			commands.DevCommand(),
			commands.BuildCommand(),
			commands.TestCommand(),
			commands.DoctorCommand(),
			commands.GenerateCommand(),
			commands.InteractiveCommand(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "Suppress non-error output",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output in JSON format (for scripting)",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "color",
				Usage: "Force colored output",
				Value: true,
			},
		},
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				util.PrintError(err.Error())
				os.Exit(1)
			}
		},
		EnableBashCompletion: true,
	}

	err := app.Run(os.Args)
	if err != nil {
		util.PrintError(err.Error())
		os.Exit(1)
	}
}
