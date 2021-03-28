package main

import (
	"log"
	"os"

	"github.com/twistedogic/task/cmd"
)

func main() {
	if err := cmd.App().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
