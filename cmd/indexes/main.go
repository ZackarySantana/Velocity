package main

import (
	"context"
	"log"

	"github.com/zackarysantana/velocity/internal/db"
)

func main() {
	ctx := context.Background()
	client, err := db.Connect(&ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Apply all indexes
	err = client.ApplyIndexes(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
