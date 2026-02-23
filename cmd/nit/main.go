package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/zGIKS/nit/internal/nit/core"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch strings.TrimSpace(strings.ToLower(os.Args[1])) {
		case "--version", "-version", "version":
			fmt.Println(version)
			return
		}
	}
	if err := core.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
