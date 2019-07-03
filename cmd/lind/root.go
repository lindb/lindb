package lind

import (
	"github.com/spf13/cobra"
)

const (
	// default config file location of LinDB
	cfgFilePath = "/etc/lindb"

	linDBText = `
    __     _             ____     ____ 
   / /    (_)   ____    / __ \   / __ )
  / /    / /   / __ \  / / / /  / __  |
 / /___ / /   / / / / / /_/ /  / /_/ / 
/_____//_/   /_/ /_/ /_____/  /_____/  

LinDB is a scalable, distributed, high performance, high availability Time Series Database, produced by Eleme-CI.
Complete documentation is available at https://github.com/eleme/lindb
`
)

// RootCmd command of cobra
var RootCmd = &cobra.Command{
	Use:   "lind",
	Short: "lind is the main command, used to control LinDB",
	Long:  linDBText,
}

func init() {
	RootCmd.AddCommand(
		envCmd,
		versionCmd,
		newStorageCmd(),
		newBrokerCmd(),
	)
}
