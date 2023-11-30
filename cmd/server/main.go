package main

import (
	"log"

	apiv1 "github.com/zackarysantana/velocity/internal/api/v1"
	"github.com/zackarysantana/velocity/internal/db"
)

func main() {
	client, err := db.Connect(nil)
	if err != nil {
		log.Fatal(err)
	}

	v1, err := apiv1.CreateV1App(*client)
	if err != nil {
		log.Fatal(err)
	}

	v1.Run(":8080")
}
