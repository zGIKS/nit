package main

import (
	"fmt"

	"nit/internal/nit/core"
)

func main() {
	if err := core.Run(); err != nil {
		fmt.Printf("error running program: %v\n", err)
	}
}
