package main

import (
	"context"
	"os"

	"github.com/zackarysantana/velocity/internal/cmd"
)

func main() {
	if err := cmd.CreateCommand().Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
