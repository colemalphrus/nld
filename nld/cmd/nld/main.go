package main

import (
	"fmt"
	"os"

	"github.com/colemalphrus/nld/internal/cli"
)

func main() {
	fmt.Println("NLD - Next-Gen Layout Document Tool")
	fmt.Println("Version: 0.1.0")

	c := cli.New()
	if err := c.Execute(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}