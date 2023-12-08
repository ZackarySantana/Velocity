package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/zackarysantana/velocity/internal/agent"
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

	stop := make(chan bool)
	wg := sync.WaitGroup{}
	ctx, err := jobs.NewCurrentContext()
	if err != nil {
		log.Fatal(err)
	}
	a := agent.NewAgent(jobs.NewVelocityJobProvider(v), &jobs.DockerJobExecutor{}, ctx, stop, &wg)

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
