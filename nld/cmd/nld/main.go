package main

import (
	"fmt"
	"os"

	"github.com/colemalphrus/nld/internal/cli"
)

func main() {
	c := cli.New()
	if err := c.Execute(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}