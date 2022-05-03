package search

import (
	"github.com/urfave/cli/v2"

	"github.com/twistedogic/task/pkg/search"
	"github.com/twistedogic/task/pkg/search/sourcegraph"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "search",
		Usage: "code search",
		Action: func(c *cli.Context) error {
			return search.New(
				c.Context,
				sourcegraph.NewWithDefault(),
			).Run(c.Args().Slice()...)
		},
	}
}
