package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/twistedogic/task/cheat"
	"github.com/twistedogic/task/docker"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "my task",
}

func init() {
	rootCmd.AddCommand(docker.RunCmd, docker.DevCmd)
	rootCmd.AddCommand(cheat.RunCmd)
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	execute()
}
