package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zackarysantana/velocity/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const databaseFile = "commands.db"

// CommandInfo represents information about a command.
type CommandInfo struct {
	ID      int
	Image   string
	Command string
	Status  *string
	Log     *string
}

func main() {
	// Create or open the database
	client, err := db.Connect(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Database opened.")

	// Create the commands table if it doesn't exist

	addExampleTasks(client.Client)
	fmt.Println("Example tasks added to the database.")
	printCommands(client.Client)
}

func printCommands(db *mongo.Client) {
	// Query all commands from the database
	cursor, err := db.Database("velocity").Collection("commands").Find(context.Background(), bson.D{{Key: "status", Value: nil}})
	if err != nil {
		log.Fatal(err)
	}

	// Print each command
	for cursor.Next(context.Background()) {
		var cmdInfo CommandInfo
		err := cursor.Decode(&cmdInfo)
		if err != nil {
			log.Fatal(err)
		}
		log := ""
		if cmdInfo.Log != nil {
			log = *cmdInfo.Log
		}

		status := ""
		if cmdInfo.Status != nil {
			status = *cmdInfo.Status
		}

		fmt.Printf("ID: %d\nImage: %s\nCommand: %s\nStatus: %s\nLog: %s\n\n",
			cmdInfo.ID, cmdInfo.Image, cmdInfo.Command, log, status)

	}
	cursor.Close(context.Background())

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

func addExampleTasks(db *mongo.Client) {
	for i := 0; i < 3; i++ {
		// Add 10 tasks every 2 seconds
		for j := 1; j <= 10; j++ {
			command := fmt.Sprintf("echo 'Task %d'", j)
			insertCommand(db, "alpine", command)
		}

		// Sleep for 2 seconds
		// time.Sleep(2 * time.Second)
	}

	fmt.Println("Example tasks added to the database.")
}

func insertCommand(db *mongo.Client, image, command string) {
	_, err := db.Database("velocity").Collection("commands").InsertOne(context.Background(), bson.D{{Key: "image", Value: image}, {Key: "command", Value: command}, {Key: "status", Value: nil}, {Key: "log", Value: nil}})
	if err != nil {
		log.Println(err)
	}
}
