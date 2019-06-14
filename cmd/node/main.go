package main

import (
	"fmt"

	"github.com/eleme/lindb/cmd/cli"
)

func main() {
	c := cli.New()

	if err := c.Run(); nil != err {
		fmt.Printf("error exit")
	}
}
