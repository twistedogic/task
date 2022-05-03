package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/twistedogic/task/cmd/dev"
	"github.com/twistedogic/task/cmd/okr"
	"github.com/twistedogic/task/cmd/run"
	"github.com/twistedogic/task/cmd/search"
)

func App() *cli.App {
	return &cli.App{
		Name:  "task",
		Usage: "Collections of utilities",
		Commands: []*cli.Command{
			dev.Command(),
			run.Command(),
			okr.Command(),
			search.Command(),
		},
	}
}
