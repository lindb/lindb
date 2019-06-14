package main

import (
	"github.com/eleme/lindb/cmd/cli"
	"fmt"
)

func main() {

	c := cli.New()

	if err := c.Run(); nil != err {
		fmt.Printf("error exit")
	}
}
