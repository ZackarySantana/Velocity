package main

import (
	"log"

	"github.com/zackarysantana/velocity/internal/operations"
)

func main() {
	if err := operations.NewCLIApp().Run(); err != nil {
		log.Fatal(err)
	}
}
