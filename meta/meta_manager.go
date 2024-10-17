package meta

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/spi/types"
)

type MetadataManager interface {
	GetStateRepo() state.Repository
	// GetMaster returns the current master info.
	GetMaster() *models.Master
	// GetBrokerNodes returns all alive broker nodes.
	GetBrokerNodes() (nodes []models.StatelessNode)
	// GetStorageNodes returns all alive storage nodes.
	GetStorageNodes() (nodes []models.StatefulNode)
	// GetDatabaseCfg returns the database config by name.
	GetDatabase(database string) (models.Database, bool)
	// GetDatabases returns current database config list.
	GetDatabases() []models.Database
	GetPartitions(database, ns, table string) (partitions map[models.InternalNode][]int, err error)
	GetTableMetadata(database, ns, table string) (*types.TableMetadata, error)
}
