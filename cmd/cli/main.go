package main

import (
	"fmt"

	"github.com/zackarysantana/velocity/internal/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Println("It seems like there was an error. If you think this is a bug, please open an issue at https://github.com/zackarysantana/velocity/issues/new/choose")
	}
}
