package lind

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// These variables are populated via the Go linker.
var (
	// release version, ldflags
	version = "unknown"
	// binary build-time, ldflags
	buildTime = "unknown"
)

func printVersion() {
	fmt.Printf("LinDB %v, BuildDate: %v\n", version, buildTime)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print environment information",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
		fmt.Printf("GOOS=%q\n", runtime.GOOS)
		fmt.Printf("GOARCH=%q\n", runtime.GOARCH)
		fmt.Printf("GOVERSION=%q\n", runtime.Version())
	},
}
