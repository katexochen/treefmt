package main

import (
	"os"

	"github.com/numtide/treefmt/v2/cmd"
)

func main() {
	// todo how are exit codes thrown by commands?
	root, _ := cmd.NewRoot()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
