package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zackarysantana/velocity/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CommandInfo represents information about a command.

type CommandInfo struct {
	ID      primitive.ObjectID `bson:"_id"`
	Command string             `bson:"command"`
	Image   string             `bson:"image"`
	Log     *string            `bson:"log"`
	Status  *string            `bson:"status"`
}

func isDockerInstalled() bool {
	return exec.Command("docker", "--version").Run() == nil
}

func main() {
	if !isDockerInstalled() {
		log.Fatal("Docker is not installed.")
	}
	client, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	commandQueue := make(chan CommandInfo)
	resultQueue := make(chan CommandInfo)
	stop := make(chan struct{})
	var wg sync.WaitGroup

	go processCommands(client, commandQueue, resultQueue, stop, &wg)
	go enqueueCommands(client, commandQueue, stop)

	go func() {
		for result := range resultQueue {
			fmt.Printf("Command %d completed with status: %s\n", result.ID, *result.Status)
		}
	}()

	fmt.Println("Press Enter to stop...")
	fmt.Scanln()

	// Stop the background process
	close(stop)
	wg.Wait()
	fmt.Println("Program terminated.")
	// printAllCommands()
}

func enqueueCommands(db *mongo.Client, commandQueue chan<- CommandInfo, stop <-chan struct{}) {
	for {
		time.Sleep(2 * time.Second)
		// select {
		// case <-stop:
		// 	return
		// default:
		// }

		cursor, err := db.Database("velocity").Collection("commands").Find(context.Background(), bson.D{{Key: "status", Value: nil}})
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = db.Database("velocity").Collection("commands").UpdateMany(context.Background(), bson.D{{Key: "status", Value: nil}}, bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: "queued"}}}})

		results := []CommandInfo{}
		// Iterate through the cursor and decode documents into your struct
		for cursor.Next(context.Background()) {
			var cmdInfo CommandInfo
			err := cursor.Decode(&cmdInfo)
			if err != nil {
				log.Fatal(err)
			}
			results = append(results, cmdInfo)
		}
		cursor.Close(context.Background())

		for _, cmdInfo := range results {
			fmt.Println("Enqueuing command: ", cmdInfo.Command)
			commandQueue <- cmdInfo
		}

		// Handle errors from cursor iteration
		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Enqueued commands.")
	}
}

func processCommands(db *mongo.Client, commandQueue <-chan CommandInfo, resultQueue chan<- CommandInfo, stop <-chan struct{}, wg *sync.WaitGroup) {
	var mu sync.Mutex
	semaphore := make(chan struct{}, 5)

	for command := range commandQueue {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(command CommandInfo) {
			defer func() {
				wg.Done()
				<-semaphore
			}()

			log.Printf("Executing command: %s", command.Command)
			l, err := executeCommand(command)
			command.Log = &l
			status := "complete"
			command.Status = &status

			fmt.Println(command.ID)

			// Update the command status
			mu.Lock()
			err = logCommand(db, command)
			if err != nil {
				log.Printf("Failed to update command: %s", err)
			}
			mu.Unlock()

			// Send the result to the result queue
			resultQueue <- command
		}(command)
	}
}

func executeCommand(command CommandInfo) (string, error) {
	cmd := exec.Command("docker", "run", "--rm", command.Image, "sh", "-c", command.Command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func logCommand(db *mongo.Client, command CommandInfo) error {
	_, err := db.Database("velocity").Collection("commands").UpdateOne(context.Background(), bson.D{{Key: "_id", Value: command.ID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: *command.Status}, {Key: "log", Value: command.Log}}}})
	return err
}
