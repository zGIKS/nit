package main

import (
	"fmt"
	"os"

	"github.com/zGIKS/nit/internal/nit/core"
)

func main() {
	if err := core.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
