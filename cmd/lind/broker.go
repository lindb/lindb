package lind

import (
	"fmt"
	_ "net/http/pprof" // for profiling

	"github.com/lindb/lindb/broker"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"

	"github.com/spf13/cobra"
)

var (
	brokerCfgPath = ""
	brokerDebug   = false
)

const (
	brokerCfgName        = "broker.toml"
	defaultBrokerCfgFile = "./" + brokerCfgName
)

// newBrokerCmd returns a new broker-cmd
func newBrokerCmd() *cobra.Command {
	brokerCmd := &cobra.Command{
		Use:     "broker",
		Aliases: []string{"bro"},
		Short:   "Run as a compute node in cluster mode",
	}
	runBrokerCmd.PersistentFlags().StringVar(&brokerCfgPath, "config", "",
		fmt.Sprintf("broker config file path, default is %s", defaultBrokerCfgFile))
	runBrokerCmd.PersistentFlags().BoolVar(&brokerDebug, "debug", false,
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
	Short: "initialize a new broker-config by steps",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := brokerCfgPath
		if len(path) == 0 {
			path = defaultBrokerCfgFile
		}
		defaultCfg := config.NewDefaultBrokerCfg()
		return fileutil.EncodeToml(path, &defaultCfg)
	},
}

// serveBroker runs the broker
func serveBroker(cmd *cobra.Command, args []string) error {
	ctx := newCtxWithSignals()

	brokerCfg := config.Broker{}
	if err := fileutil.LoadConfig(brokerCfgPath, defaultBrokerCfgFile, &brokerCfg); err != nil {
		return fmt.Errorf("decode config file error:%s", err)
	}

	// start broker server
	brokerRuntime := broker.NewBrokerRuntime(brokerCfg)
	if err := run(ctx, brokerRuntime); err != nil {
		return err
	}
	return nil
}
