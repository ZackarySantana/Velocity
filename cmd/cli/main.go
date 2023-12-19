package main

import (
	"github.com/zackarysantana/velocity/internal/operations"
)

func main() {
	if err := operations.NewCLIApp().Run(); err != nil {
		// Portentially log error on server?
	}
}
