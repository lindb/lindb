package query

import (
	"errors"
)

var errNoAvailableStorageNode = errors.New("no available storage node for server")
var errDatabaseNotExist = errors.New("database not exist")
