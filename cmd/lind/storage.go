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

// newStorageCmd returns a new storage-cmd
func newStorageCmd() *cobra.Command {
	storageCmd := &cobra.Command{
		Use:     "storage",
		Aliases: []string{"sto", "stor"},
		Short:   "The storage layer of LinDB",
	}
	runStorageCmd.PersistentFlags().StringVar(&storageCfgPath, "config", "",
		fmt.Sprintf("storage config file path, default is %s", storage.DefaultStorageCfgFile))
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
	RunE: func(cmd *cobra.Command, args []string) error {
		path := storageCfgPath
		if len(path) == 0 {
			path = storage.DefaultStorageCfgFile
		}
		defaultCfg := config.NewDefaultStorageCfg()
		return fileutil.EncodeToml(path, &defaultCfg)
	},
}

func serveStorage(cmd *cobra.Command, args []string) error {
	ctx := newCtxWithSignals()

	// start storage server
	storageRuntime := storage.NewStorageRuntime(storageCfgPath)
	if err := storageRuntime.Run(); err != nil {
		return fmt.Errorf("run storage server error:%s", err)
	}

	// waiting system exit signal
	<-ctx.Done()

	// stop storage server
	if err := storageRuntime.Stop(); err != nil {
		return fmt.Errorf("stop storage server error:%s", err)
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
