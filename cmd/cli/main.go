package main

import (
	"context"
	"log"
	"os"

	"github.com/zackarysantana/velocity/internal/cmd"
)

func main() {
	if err := cmd.CreateCommand().Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
