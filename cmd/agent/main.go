package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/zackarysantana/velocity/internal/agent"
	"github.com/zackarysantana/velocity/internal/jobs"
	"github.com/zackarysantana/velocity/src/clients"
)

func main() {
	v := clients.NewVelocityClientV1WithAPIKey("http://localhost:8080", "YOUR_API_KEY")

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
