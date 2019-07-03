package lind

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // for profiling
	"os"

	"github.com/eleme/lindb/broker/rest"
	"github.com/eleme/lindb/config"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var (
	brokerCfgPath = ""
	brokerDebug   = false
)

const (
	brokerCfgName        = "broker.toml"
	defaultBrokerCfgFile = cfgFilePath + "/" + brokerCfgName
)

// newBrokerCmd returns a new broker-cmd
func newBrokerCmd() *cobra.Command {
	brokerCmd := &cobra.Command{
		Use:     "broker",
		Aliases: []string{"bro"},
		Short:   "The compute layer of LinDB",
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
	Use:   "initialize-config",
	Short: "initialize a new broker-config by steps",
	Run: func(cmd *cobra.Command, args []string) {
		// todo: @codingcrush
	},
}

// serveBroker runs the broker
func serveBroker(cmd *cobra.Command, args []string) error {
	ctx := newCtxWithSignals()
	go func() {
		<-ctx.Done()
		os.Exit(0)
	}()

	if brokerCfgPath == "" {
		brokerCfgPath = defaultBrokerCfgFile
	}
	if _, err := os.Stat(brokerCfgPath); err != nil {
		return fmt.Errorf("config file doesn't exist, see how to initialize the config by `lind broker -h`")
	}
	fmt.Printf("load config file: %v successfully\n", brokerCfgPath)

	brokerConfig := config.BrokerConfig{}
	if _, err := toml.DecodeFile(brokerCfgPath, &brokerConfig); err != nil {
		return err
	}
	fmt.Printf("HTTP server listening on: %d\n", brokerConfig.HTTP.Port)

	router := rest.NewRouter(&brokerConfig)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", brokerConfig.HTTP.Port), router); err != nil {
		return err
	}
	return nil
}
