package main

import (
	"fmt"
	"os"

	"github.com/bensallen/rbd/internal/cli/root"
	"github.com/bensallen/rbd/pkg/boot"
)

func main() {
	boot.PIDInit()
	if err := root.Run(os.Args[1:]); err != nil {
		fmt.Printf("Error: %v\n\n", err)
		os.Exit(1)
	}
}
