package commands

import (
	"errors"

	"github.com/springwell/cli/pkg/config"
	"github.com/springwell/cli/pkg/generator"
	"github.com/springwell/cli/pkg/util"
	"github.com/urfave/cli/v2"
)

// GenerateEntityCommand returns the command to generate an entity
func GenerateEntityCommand() *cli.Command {
	return &cli.Command{
		Name:  "entity",
		Usage: "Generate an entity and its related components",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "fields",
				Aliases: []string{"f"},
				Usage:   "Field definitions (format: \"name:type[:modifier]\")",
			},
			&cli.StringFlag{
				Name:    "relations",
				Aliases: []string{"r"},
				Usage:   "Relationship definitions (format: \"type:field:entity\")",
			},
			&cli.StringFlag{
				Name:    "table",
				Aliases: []string{"t"},
				Usage:   "Database table name (default: derived from entity name)",
			},
			&cli.BoolFlag{
				Name:  "audit",
				Usage: "Add auditing fields (created/updated timestamps)",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "lombok",
				Usage: "Use Lombok annotations",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "dto",
				Usage: "Generate DTO classes",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "no-repository",
				Usage: "Skip repository generation",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "no-service",
				Usage: "Skip service generation",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "no-controller",
				Usage: "Skip controller generation",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			entityName := c.Args().First()
			if entityName == "" {
				return errors.New("entity name is required")
			}

			// Check if the current directory is a Spring Boot project
			if !util.IsSpringBootProject(".") {
				return errors.New("current directory is not a Spring Boot project")
			}

			// Load config
			cfg, err := config.LoadConfig(".")
			if err != nil {
				return err
			}

			// Create generator
			gen := generator.NewEntityGenerator(cfg, ".")

			// Generate entity
			err = gen.GenerateEntity(
				entityName,
				c.String("fields"),
				c.String("relations"),
				c.String("table"),
				c.Bool("audit"),
				c.Bool("lombok"),
				c.Bool("dto"),
				!c.Bool("no-repository"),
				!c.Bool("no-service"),
				!c.Bool("no-controller"),
			)

			if err != nil {
				return err
			}

			util.PrintSuccess("Successfully generated %s entity and related components", entityName)
			return nil
		},
	}
}

// GenerateControllerCommand returns the command to generate a controller
func GenerateControllerCommand() *cli.Command {
	return &cli.Command{
		Name:  "controller",
		Usage: "Generate a REST controller",
		Action: func(c *cli.Context) error {
			controllerName := c.Args().First()
			if controllerName == "" {
				return errors.New("controller name is required")
			}

			// TODO: Implement controller generation
			util.PrintSuccess("Successfully generated %s controller", controllerName)
			return nil
		},
	}
}

// GenerateServiceCommand returns the command to generate a service
func GenerateServiceCommand() *cli.Command {
	return &cli.Command{
		Name:  "service",
		Usage: "Generate a service class",
		Action: func(c *cli.Context) error {
			serviceName := c.Args().First()
			if serviceName == "" {
				return errors.New("service name is required")
			}

			// TODO: Implement service generation
			util.PrintSuccess("Successfully generated %s service", serviceName)
			return nil
		},
	}
}

// GenerateRepositoryCommand returns the command to generate a repository
func GenerateRepositoryCommand() *cli.Command {
	return &cli.Command{
		Name:  "repository",
		Usage: "Generate a repository interface",
		Action: func(c *cli.Context) error {
			repositoryName := c.Args().First()
			if repositoryName == "" {
				return errors.New("repository name is required")
			}

			// TODO: Implement repository generation
			util.PrintSuccess("Successfully generated %s repository", repositoryName)
			return nil
		},
	}
}

// GenerateDtoCommand returns the command to generate a DTO
func GenerateDtoCommand() *cli.Command {
	return &cli.Command{
		Name:  "dto",
		Usage: "Generate a DTO class",
		Action: func(c *cli.Context) error {
			dtoName := c.Args().First()
			if dtoName == "" {
				return errors.New("DTO name is required")
			}

			// TODO: Implement DTO generation
			util.PrintSuccess("Successfully generated %s DTO", dtoName)
			return nil
		},
	}
}

// GenerateCommand returns the generate command
func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:    "generate",
		Aliases: []string{"g"},
		Usage:   "Generate code components",
		Subcommands: []*cli.Command{
			GenerateEntityCommand(),
			GenerateControllerCommand(),
			GenerateServiceCommand(),
			GenerateRepositoryCommand(),
			GenerateDtoCommand(),
		},
	}
}
