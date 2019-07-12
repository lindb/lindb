package lind

import (
	"fmt"
	_ "net/http/pprof" // for profiling

	"github.com/eleme/lindb/broker"
	"github.com/eleme/lindb/config"
	"github.com/eleme/lindb/pkg/util"

	"github.com/spf13/cobra"
)

var (
	brokerCfgPath = ""
	brokerDebug   = false
)

// newBrokerCmd returns a new broker-cmd
func newBrokerCmd() *cobra.Command {
	brokerCmd := &cobra.Command{
		Use:     "broker",
		Aliases: []string{"bro"},
		Short:   "The compute layer of LinDB",
	}
	runBrokerCmd.PersistentFlags().StringVar(&brokerCfgPath, "config", "",
		fmt.Sprintf("broker config file path, default is %s", broker.DefaultBrokerCfgFile))
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
	Use:   "initialize-config",
	Short: "initialize a new broker-config by steps",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := brokerCfgPath
		if len(path) == 0 {
			path = broker.DefaultBrokerCfgFile
		}
		defaultCfg := config.NewDefaultBrokerCfg()
		return util.EncodeToml(path, &defaultCfg)
	},
}

// serveBroker runs the broker
func serveBroker(cmd *cobra.Command, args []string) error {
	ctx := newCtxWithSignals()

	// start broker server
	brokerRuntime := broker.NewBrokerRuntime(brokerCfgPath)
	if err := brokerRuntime.Run(); err != nil {
		return fmt.Errorf("run broker server error:%s", err)
	}

	// waiting system exit signal
	<-ctx.Done()

	// stop broker server
	if err := brokerRuntime.Stop(); err != nil {
		return fmt.Errorf("stop broker server error:%s", err)
	}

	return nil
}
