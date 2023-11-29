package main

// Simple web server that serves the word HELLOWORLD on port 8080

import (
	"net/http"

	"github.com/zackarysantana/repo1/pkg1"
)

func main() {
	http.HandleFunc("/", pkg1.HelloWorld)
	http.ListenAndServe(":8080", nil)
}
