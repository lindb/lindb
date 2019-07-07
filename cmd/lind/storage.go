package lind

import (
	"fmt"
	_ "net/http/pprof" // for profiling
	"os"

	"github.com/eleme/lindb/config"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/query"
	"github.com/eleme/lindb/storage"
	"github.com/eleme/lindb/tsdb"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	storageCfgName        = "storage.toml"
	defaultStorageCfgFile = cfgFilePath + "/" + storageCfgName
)

var (
	storageCfgPath = ""
	storageDebug   = false
)

// newStorageCmd returns a new storage-cmd
func newStorageCmd() *cobra.Command {
	storageCmd := &cobra.Command{
		Use:     "storage",
		Aliases: []string{"sto", "stor"},
		Short:   "The storage layer of LinDB",
	}
	runStorageCmd.PersistentFlags().StringVar(&storageCfgPath, "config", "",
		fmt.Sprintf("storage config file path, default is %s", defaultStorageCfgFile))
	runStorageCmd.PersistentFlags().BoolVar(&storageDebug, "debug", false,
		"profiling Go programs with pprof")

	storageCmd.AddCommand(
		runStorageCmd,
		initializeStorageConfigCmd,
		databaseCmd,
	)
	return storageCmd
}

var runStorageCmd = &cobra.Command{
	Use:   "run",
	Short: "starts the storage",
	RunE:  serveStorage,
}

var initializeStorageConfigCmd = &cobra.Command{
	Use:   "initialize-config",
	Short: "initialize a new storage-config by steps",
	Run: func(cmd *cobra.Command, args []string) {
		// todo: @codingcrush
	},
}

func serveStorage(cmd *cobra.Command, args []string) error {
	log := logger.GetLogger()
	ctx := newCtxWithSignals()

	if storageCfgPath == "" {
		storageCfgPath = defaultStorageCfgFile
	}
	if _, err := os.Stat(storageCfgPath); err != nil {
		return fmt.Errorf("config file doesn't exist, see how to initialize the config by `lind storage -h`")
	}
	fmt.Printf("load config file: %v successfully\n", storageCfgPath)

	storageConfig := config.StorageConfig{}
	if _, err := toml.DecodeFile(storageCfgPath, &storageConfig); err != nil {
		return err
	}
	// start the repository server
	storageServer := storage.New(ctx, &storageConfig)
	if err := storageServer.Start(); err != nil {
		log.Error("storage start failed", zap.Error(err))
		return err
	}

	// create a new engine
	engine, err := tsdb.NewEngine(storageConfig.Name, storageConfig.Path)
	if err != nil {
		return err
	}
	// todo: fix this
	_ = query.NewTSDBExecutor(engine, nil, nil)

	<-ctx.Done()
	return engine.Close()
}

// databaseCmd provides the ability to control the database of storage
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Control the database of LinDB",
	Run: func(cmd *cobra.Command, args []string) {
		// todo: @codingcrush
	},
}
