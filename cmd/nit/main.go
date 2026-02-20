package main

import (
	"fmt"

	"nit/internal/nit"
)

func main() {
	if err := nit.Run(); err != nil {
		fmt.Printf("error running program: %v\n", err)
	}
}
