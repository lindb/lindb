package meta

import (
	"github.com/lindb/lindb/models"
)

type MetadataManager interface {
	// GetMaster returns the current master info.
	GetMaster() *models.Master
	// GetDatabaseCfg returns the database config by name.
	GetDatabase(database string) (models.Database, bool)
	// GetDatabases returns current database config list.
	GetDatabases() []models.Database
	GetPartitions(database, ns, table string) (partitions map[models.InternalNode][]int, err error)
}
