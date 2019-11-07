package lind

import (
	"fmt"

	"github.com/lindb/lindb/broker"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"

	"github.com/spf13/cobra"
)

const (
	brokerCfgName        = "broker.toml"
	defaultBrokerCfgFile = "./" + brokerCfgName
)

// newBrokerCmd returns a new broker-cmd
func newBrokerCmd() *cobra.Command {
	brokerCmd := &cobra.Command{
		Use:   "broker",
		Short: "Run as a compute node with cluster mode enabled",
	}
	runBrokerCmd.PersistentFlags().StringVar(&cfg, "config", "",
		fmt.Sprintf("broker config file path, default is %s", defaultBrokerCfgFile))
	runBrokerCmd.PersistentFlags().BoolVar(&debug, "debug", false,
		"profiling Go programs with pprof")
	brokerCmd.AddCommand(
		runBrokerCmd,
		initializeBrokerConfigCmd,
	)
	return brokerCmd
}

var runBrokerCmd = &cobra.Command{
	Use:   "run",
	Short: "starts the broker",
	RunE:  serveBroker,
}

// initialize config for broker
var initializeBrokerConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "create a new default broker-config",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cfg
		if len(path) == 0 {
			path = defaultBrokerCfgFile
		}
		if err := checkExistenceOf(path); err != nil {
			return err
		}
		defaultCfg := config.NewDefaultBrokerCfg()
		return fileutil.EncodeToml(path, &defaultCfg)
	},
}

// serveBroker runs the broker
func serveBroker(cmd *cobra.Command, args []string) error {
	ctx := newCtxWithSignals()

	brokerCfg := config.Broker{}
	if err := fileutil.LoadConfig(cfg, defaultBrokerCfgFile, &brokerCfg); err != nil {
		return fmt.Errorf("decode config file error: %s", err)
	}
	if err := logger.InitLogger(brokerCfg.Logging); err != nil {
		return fmt.Errorf("init logger error: %s", err)
	}

	// start broker server
	brokerRuntime := broker.NewBrokerRuntime(getVersion(), brokerCfg)
	if err := run(ctx, brokerRuntime); err != nil {
		return err
	}
	return nil
}
