package main

import (
	"github.com/zackarysantana/velocity/cmd/analysis/vets"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(vets.CommandNoFmt)
}
