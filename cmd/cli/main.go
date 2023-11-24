package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	operation "github.com/zackarysantana/velocity/internal/operations"
)

func main() {
	app := &cli.App{
		Name:     "velocity",
		Usage:    "manage, run, and report on tests quickly",
		Commands: append([]cli.Command{}, operation.Run...),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
