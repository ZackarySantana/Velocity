package main

import (
	"context"
	"fmt"

	"github.com/zackarysantana/velocity/internal/db"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	client, err := db.Connect()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}
