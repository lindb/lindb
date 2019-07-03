package lind

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// These variables are populated via the Go linker.
var (
	// release version, ldflags
	version string
	// binary build-time, ldflags
	buildTime string
)

const (
	defaultVersion = "alpha"
)

func printVersion() {
	var releaseVersion = defaultVersion
	if version != "" {
		releaseVersion = version
	}
	fmt.Printf("LinDB %v, BuildDate: %v\n", releaseVersion, buildTime)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of LinDB",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print environment info of LinDB",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
		fmt.Printf("GOOS=%q\n", runtime.GOOS)
		fmt.Printf("GOARCH=%q\n", runtime.GOARCH)
		fmt.Printf("GOVERSION=%q\n", runtime.Version())
	},
}
