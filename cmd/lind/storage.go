package lind

import (
	"fmt"
	_ "net/http/pprof" // for profiling

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/storage"

	"github.com/spf13/cobra"
)

var (
	storageCfgPath = ""
	storageDebug   = false
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
		Use:     "storage",
		Aliases: []string{"sto", "stor"},
		Short:   "Run as a storage node in cluster mode",
	}
	runStorageCmd.PersistentFlags().BoolVar(&storageDebug, "debug", false,
		"profiling Go programs with pprof")
	runStorageCmd.PersistentFlags().StringVar(&storageCfgPath, "config", "",
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
	Short: "initialize a new storage-config by steps",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := storageCfgPath
		if len(path) == 0 {
			path = defaultStorageCfgFile
		}
		defaultCfg := config.NewDefaultStorageCfg()
		return fileutil.EncodeToml(path, &defaultCfg)
	},
}

func serveStorage(cmd *cobra.Command, args []string) error {
	ctx := newCtxWithSignals()

	storageCfg := config.Storage{}
	if err := fileutil.LoadConfig(storageCfgPath, defaultStorageCfgFile, &storageCfg); err != nil {
		return fmt.Errorf("decode config file error:%s", err)
	}
	// start storage server
	storageRuntime := storage.NewStorageRuntime(storageCfg)
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
