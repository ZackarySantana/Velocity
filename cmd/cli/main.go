package main

import (
	"context"
	"os"

	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/internal/cmd"
)

func main() {
	oops.StackTraceMaxDepth = 0

	if err := cmd.CreateCommand().Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
