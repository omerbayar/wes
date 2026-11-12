package main

import (
	"os"

	"github.com/omerbayar/wes/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
