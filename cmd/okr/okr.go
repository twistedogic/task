package okr

import (
	"github.com/urfave/cli/v2"

	"github.com/twistedogic/task/pkg/okr"
)

func AddCmd() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add Objective Key Result",
		Action: func(c *cli.Context) error {
			okr.Add()
			return nil
		},
	}
}

func EditCmd() *cli.Command {
	return &cli.Command{
		Name:  "edit",
		Usage: "Edit Objective Key Result",
		Action: func(c *cli.Context) error {
			okr.Edit()
			return nil
		},
	}
}

func UpdateCmd() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Update Key Result Progress",
		Action: func(c *cli.Context) error {
			okr.Update()
			return nil
		},
	}
}

func Command() *cli.Command {
	return &cli.Command{
		Name:  "okr",
		Usage: "Objective Key Result",
		Subcommands: []*cli.Command{
			AddCmd(),
			EditCmd(),
			UpdateCmd(),
		},
		Action: func(c *cli.Context) error {
			okr.List()
			return nil
		},
	}
}
