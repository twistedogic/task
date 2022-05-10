package okr

import (
	"github.com/urfave/cli/v2"

	"github.com/twistedogic/task/pkg/okr"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "okr",
		Usage: "Objective Key Result",
		Action: func(c *cli.Context) error {
			app, err := okr.New(c.Context)
			if err != nil {
				return err
			}
			return app.Run()
		},
	}
}
