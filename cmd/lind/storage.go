package lind

import (
	"fmt"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/storage"

	"github.com/spf13/cobra"
)

const (
	storageCfgName        = "storage.toml"
	defaultStorageCfgFile = "./" + storageCfgName
)

var runStorageCmd = &cobra.Command{
	Use:   "run",
	Short: "starts the storage",
	RunE:  serveStorage,
}

// newStorageCmd returns a new storage-cmd
func newStorageCmd() *cobra.Command {
	storageCmd := &cobra.Command{
		Use:   "storage",
		Short: "Run as a storage node with cluster mode enabled",
	}
	runStorageCmd.PersistentFlags().BoolVar(&debug, "debug", false,
		"profiling Go programs with pprof")
	runStorageCmd.PersistentFlags().StringVar(&cfg, "config", "",
		fmt.Sprintf("storage config file path, default is %s", defaultStorageCfgFile))

	storageCmd.AddCommand(
		runStorageCmd,
		initializeStorageConfigCmd,
		databaseCmd,
	)
	return storageCmd
}

var initializeStorageConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "create a new default storage-config",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cfg
		if len(path) == 0 {
			path = defaultStorageCfgFile
		}
		if err := checkExistenceOf(path); err != nil {
			return err
		}
		defaultCfg := config.NewDefaultStorageCfg()
		return fileutil.EncodeToml(path, &defaultCfg)
	},
}

func serveStorage(cmd *cobra.Command, args []string) error {
	ctx := newCtxWithSignals()

	storageCfg := config.Storage{}
	if err := fileutil.LoadConfig(cfg, defaultStorageCfgFile, &storageCfg); err != nil {
		return fmt.Errorf("decode config file error: %s", err)
	}
	if err := logger.InitLogger(storageCfg.Logging); err != nil {
		return fmt.Errorf("init logger error: %s", err)
	}

	// start storage server
	storageRuntime := storage.NewStorageRuntime(getVersion(), storageCfg)
	if err := run(ctx, storageRuntime); err != nil {
		return err
	}
	return nil
}

// databaseCmd provides the ability to control the database of storage
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Control the database of LinDB",
	Run: func(cmd *cobra.Command, args []string) {
		// todo: @codingcrush
	},
}
