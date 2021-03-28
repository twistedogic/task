package run

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/twistedogic/task/pkg/docker"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "run container as command",
		Action: func(c *cli.Context) error {
			out, err := docker.RunTask(c.Args().Slice()...)
			if err != nil {
				return err
			}
			fmt.Println(string(out))
			return nil
		},
	}
}
