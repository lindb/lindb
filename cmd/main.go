package main

import (
	"fmt"
	"os"

	"github.com/lindb/lindb/cmd/lind"
)

func main() {
	if err := lind.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
