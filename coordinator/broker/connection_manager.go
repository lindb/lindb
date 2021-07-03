package broker

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
)

var connectionManagerLogger = logger.GetLogger("coordinator", "ConnectionManager")

// connectionManger manages the rpc connections
// not thread-safe
type connectionManager struct {
	RoleFrom          string
	RoleTo            string
	connections       map[string]struct{}
	taskClientFactory rpc.TaskClientFactory
}

func (manager *connectionManager) createConnection(target models.Node) {
	if err := manager.taskClientFactory.CreateTaskClient(target); err == nil {
		connectionManagerLogger.Info("established connection successfully",
			logger.String("target", target.Indicator()),
			logger.String("from", manager.RoleFrom),
			logger.String("to", manager.RoleTo),
		)
		manager.connections[target.Indicator()] = struct{}{}
	} else {
		connectionManagerLogger.Error("failed to establish connection",
			logger.String("target", target.Indicator()),
			logger.String("from", manager.RoleFrom),
			logger.String("to", manager.RoleTo),
			logger.Error(err),
		)
		delete(manager.connections, target.Indicator())
	}
}

func (manager *connectionManager) closeConnection(target string) {
	closed, err := manager.taskClientFactory.CloseTaskClient(target)
	delete(manager.connections, target)

	if closed {
		if err == nil {
			connectionManagerLogger.Info("closed connection successfully",
				logger.String("target", target),
				logger.String("from", manager.RoleFrom),
				logger.String("to", manager.RoleTo),
			)
		} else {
			connectionManagerLogger.Error("failed to close connection",
				logger.String("target", target),
				logger.String("from", manager.RoleFrom),
				logger.String("to", manager.RoleTo),
				logger.Error(err),
			)
		}
	} else {
		connectionManagerLogger.Debug("unable to close a non-existent connection",
			logger.String("target", target),
			logger.String("from", manager.RoleFrom),
			logger.String("to", manager.RoleTo),
		)
	}
}

func (manager *connectionManager) closeAll() {
	for target := range manager.connections {
		manager.closeConnection(target)
	}
}

func (manager *connectionManager) closeInactiveNodeConnections(activeNodes []string) {
	activeNodesSet := make(map[string]struct{})
	for _, node := range activeNodes {
		activeNodesSet[node] = struct{}{}
	}
	for target := range manager.connections {
		if _, exist := activeNodesSet[target]; !exist {
			manager.closeConnection(target)
		}
	}
}
