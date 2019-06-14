package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/eleme/lindb/broker"
	"github.com/eleme/lindb/broker/rest"
	"github.com/eleme/lindb/pkg/config"
	"github.com/eleme/lindb/pkg/logger"

	"go.uber.org/zap"
)

const unknown = "unknown"

// These variables are populated via the Go linker.
var (
	version    string
	commit     string
	configFile string
	help       bool
)

func init() {
	if version == "" {
		version = unknown
	}

	if commit == "" {
		commit = unknown
	}

	flag.BoolVar(&help, "help", false, "help")
	flag.StringVar(&configFile, "config", "/etc/lindb/broker.toml", "config file path")
}

func main() {
	flag.Parse()

	if help {
		// display help
		usage()
		os.Exit(0)
	}

	log := logger.GetLogger()

	log.Info("load config file", zap.String("path", configFile))

	brokerConfig := &broker.Config{}
	config.Parse(configFile, brokerConfig)
	log.Info("start http server", zap.Any("port", brokerConfig.HTTP.Port))

	router := rest.NewRouter(brokerConfig)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", brokerConfig.HTTP.Port), router); err != nil {
		log.Error("start http server error", zap.Error(err))
		os.Exit(0)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `lindb broker version: %s

Options:
`, version)
	flag.PrintDefaults()
}
