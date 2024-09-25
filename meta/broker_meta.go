package meta

import (
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/models"
)

type brokerMetadataManager struct {
	brokerStateMgr broker.StateManager
	masterStateMgr master.StateManager
}

func NewBrokerMetadataManager(brokerStateMgr broker.StateManager, masterStateMgr master.StateManager) MetadataManager {
	return &brokerMetadataManager{
		brokerStateMgr: brokerStateMgr,
		masterStateMgr: masterStateMgr,
	}
}

func (m *brokerMetadataManager) GetDatabase(database string) (models.Database, bool) {
	return m.brokerStateMgr.GetDatabase(database)
}

func (m *brokerMetadataManager) GetDatabases() []models.Database {
	return m.brokerStateMgr.GetDatabases()
}

func (m *brokerMetadataManager) GetPartitions(database, ns, table string) (map[models.InternalNode][]int, error) {
	if database == constants.InformationSchema {
		var partitions map[models.InternalNode][]int
		currentNode := m.brokerStateMgr.GetCurrentNode()
		switch table {
		case constants.TableSchemata:
			partitions = map[models.InternalNode][]int{
				{IP: currentNode.HostIP, Port: currentNode.GRPCPort}: {},
			}
		}
		return partitions, nil
	}

	// find tabel metadata from partitions
	return m.brokerStateMgr.GetPartitions(database)
}
