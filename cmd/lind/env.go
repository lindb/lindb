package lind

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// These variables are populated via the Go linker.
var (
	// release version, ldflags
	version = ""
	// binary build-time, ldflags
	buildTime = "unknown"
	// debug mode
	debug = false
	// cfg path
	cfg = ""
)

const defaultVersion = "0.0.0"

func getVersion() string {
	if version == "" {
		return defaultVersion
	}
	return version
}

func printVersion() {
	fmt.Printf("LinDB: %v, BuildDate: %v\n", getVersion(), buildTime)
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
