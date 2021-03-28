package dev

import (
	"github.com/urfave/cli/v2"

	"github.com/twistedogic/task/pkg/docker"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "dev",
		Usage: "run container for development",
		Action: func(c *cli.Context) error {
			return docker.RunDevEnv(c.Args().First())
		},
	}
}
