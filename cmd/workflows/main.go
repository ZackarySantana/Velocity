package main

import (
	"fmt"

	"github.com/zackarysantana/velocity/src/config"
)

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

func main() {
	config, err := config.ReadConfigFromFile("velocity.yml")
	if err != nil {
		panic(err)
	}

	// Loop through tests
	for testName, test := range config.Tests {
		fmt.Println(testName)
		fmt.Println(test)
	}
}
