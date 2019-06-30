package main

import (
	"fmt"

	"github.com/eleme/lindb/cmd/cli"
)

const lindbText = `
 _      _____  _   _ ______ ______ 
| |    |_   _|| \ | ||  _  \| ___ \
| |      | |  |  \| || | | || |_/ /
| |      | |  | .   || | | || ___ \
| |____ _| |_ | |\  || |/ / | |_/ /
\_____/ \___/ \_| \_/|___/  \____/
`

func main() {
	c := cli.New()
	fmt.Print(lindbText)
	if err := c.Run(); nil != err {
		fmt.Printf("error exit")
	}
}
