package main

import (
	"os"

	"github.com/dnagent-app/dna/internal/cli"
)

func main() {
	if err := cli.Run(os.Args); err != nil {
		cli.Fatal(err)
	}
}
