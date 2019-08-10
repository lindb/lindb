package task

import "fmt"

const (
	version = "v1"
	//Notice: magic number, see also: --max-txn-ops in etcd
	maxTasksLimit = 127
)

var (
	taskCoordinatorKey = fmt.Sprintf("/task-coordinator/%s", version)
)
