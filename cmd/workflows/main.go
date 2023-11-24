package main

import (
	"fmt"
	"os"

	"github.com/zackarysantana/velocity/src/config"
)

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: filename")
		return
	}
	filename := os.Args[1]
	config, err := config.ReadConfigFromFile(filename)
	if err != nil {
		panic(err)
	}

	// Loop through tests
	for testName, test := range config.Tests {
		fmt.Println(testName)
		fmt.Println(test)
	}
}
