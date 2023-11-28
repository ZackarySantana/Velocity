package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/zackarysantana/velocity/internal/agent"
	"github.com/zackarysantana/velocity/internal/db"
	"github.com/zackarysantana/velocity/internal/jobs"
)

func main() {
	client, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	stop := make(chan bool)
	wg := sync.WaitGroup{}
	ctx := jobs.NewContext("https://github.com/zackarysantana/velocity.git", "c8dc99dfc0b62842b0a524fe34112c3df27f7e86")
	a := agent.NewAgent(jobs.NewMongoDBJobProvider(*client), &jobs.DockerJobExecutor{}, ctx, stop, &wg)

	err = a.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Press Enter to stop...")
	fmt.Scanln()

	// Stop the background process
	close(stop)
	wg.Wait()
	fmt.Println("Program terminated.")
}
