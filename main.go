package main

import (
	"os"

	"github.com/syauqeesy/liveness-detection/foundation"
)

func main() {
	if err := foundation.Boot(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
