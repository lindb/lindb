package lind

import (
	"fmt"
	_ "net/http/pprof" // for profiling

	"github.com/spf13/cobra"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/standalone"
)

const (
	standaloneCfgName = "standalone.toml"
	// DefaultStandaloneCfgFile defines default config file path for standalone mode
	defaultStandaloneCfgFile = "./" + standaloneCfgName
)

// newStandaloneCmd returns a new standalone-cmd
func newStandaloneCmd() *cobra.Command {
	standaloneCmd := &cobra.Command{
		Use:   "standalone",
		Short: "Run as the standalone mode(embed broker/storage/etcd)",
	}

	standaloneCmd.AddCommand(
		runStandaloneCmd,
		initializeStandaloneConfigCmd,
	)

	runStandaloneCmd.PersistentFlags().BoolVar(&debug, "debug", false,
		"profiling Go programs with pprof")
	runStandaloneCmd.PersistentFlags().StringVar(&cfg, "config", "",
		fmt.Sprintf("config file path for standalone mode, default is %s", defaultStandaloneCfgFile))

	return standaloneCmd
}

var runStandaloneCmd = &cobra.Command{
	Use:   "run",
	Short: "run as the standalone mode",
	RunE:  serveStandalone,
}

// initializeStandaloneConfigCmd initializes config for standalone mode
var initializeStandaloneConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "create a new default standalone-config",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cfg
		if len(path) == 0 {
			path = defaultStandaloneCfgFile
		}
		if err := checkExistenceOf(path); err != nil {
			return err
		}
		defaultCfg := config.NewDefaultStandaloneCfg()
		return fileutil.EncodeToml(path, &defaultCfg)
	},
}

// serveStandalone runs the cluster as standalone mode
func serveStandalone(cmd *cobra.Command, args []string) error {
	ctx := newCtxWithSignals()

	standaloneCfg := config.Standalone{}
	if err := fileutil.LoadConfig(cfg, defaultStandaloneCfgFile, &standaloneCfg); err != nil {
		return fmt.Errorf("decode config file error: %s", err)
	}
	if err := logger.InitLogger(standaloneCfg.Logging); err != nil {
		return fmt.Errorf("init logger error: %s", err)
	}

	// run cluster as standalone mode
	runtime := standalone.NewStandaloneRuntime(standaloneCfg)
	if err := run(ctx, runtime); err != nil {
		return err
	}
	return nil
}
