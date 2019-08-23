package lind

import (
	"fmt"
	"os"

	"github.com/lindb/lindb/pkg/logger"

	"github.com/spf13/cobra"
)

const linDBLogo = `
██╗     ██╗███╗   ██╗██████╗ ██████╗ 
██║     ██║████╗  ██║██╔══██╗██╔══██╗
██║     ██║██╔██╗ ██║██║  ██║██████╔╝
██║     ██║██║╚██╗██║██║  ██║██╔══██╗
███████╗██║██║ ╚████║██████╔╝██████╔╝
╚══════╝╚═╝╚═╝  ╚═══╝╚═════╝ ╚═════╝ 
`

func printLogoWhenIsTty() {
	if logger.IsTerminal(os.Stdout) {
		fmt.Fprintf(os.Stdout, logger.Cyan.Add(linDBLogo))
		fmt.Fprintf(os.Stdout, logger.Green.Add(" ::  LinDB  :: ")+
			fmt.Sprintf("%22s", fmt.Sprintf("(v%s Release)", getVersion())))
		fmt.Fprintf(os.Stdout, "\n\n")
	}
}

const (
	linDBText = `
LinDB is a scalable, high performance, high availability, distributed time series database.
Complete documentation is available at https://lindb.io
`
)

// RootCmd command of cobra
var RootCmd = &cobra.Command{
	Use:   "lind",
	Short: "lind is the main command, used to control LinDB",
	Long:  linDBLogo + linDBText,
}

func init() {
	RootCmd.AddCommand(
		envCmd,
		versionCmd,
		newStorageCmd(),
		newBrokerCmd(),
		newStandaloneCmd(),
	)
}
