package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/zackarysantana/velocity/internal/agent"
	"github.com/zackarysantana/velocity/internal/db"
	"github.com/zackarysantana/velocity/internal/jobs"
	"github.com/zackarysantana/velocity/src/clients"
	"github.com/zackarysantana/velocity/src/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	v, err := clients.NewVelocityClientV1FromConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	var provider jobs.JobProvider
	if _, ok := os.LookupEnv("MONGODB_AGENT"); ok {
		client, err := db.Connect(nil)
		if err != nil {
			log.Fatal(err)
		}
		provider = jobs.NewMongoDBJobProvider(*client)
	} else {
		provider = jobs.NewVelocityJobProvider(v)
	}

	stop := make(chan bool)
	wg := sync.WaitGroup{}
	a := agent.NewAgent(provider, &jobs.DockerJobExecutor{}, stop, &wg)

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
