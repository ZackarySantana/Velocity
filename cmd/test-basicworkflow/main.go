package main

import (
	"fmt"
	"log"

	"github.com/zackarysantana/velocity/internal/workflows"
	"github.com/zackarysantana/velocity/src/config"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	w, err := workflows.GetWorkflow(*c, "Please select a workflow: ")
	if err != nil {
		log.Fatal(err)
	}
	results, err := workflows.RunSyncWorkflow(*c, w)
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		if result.Success != nil {
			fmt.Println("'" + result.Job.Image + "' ran '" + result.Job.Command + "'")
			fmt.Println(result.Success.Logs)
			fmt.Println()
		}
	}
}
