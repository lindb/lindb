package models

import (
	"github.com/lindb/lindb/pkg/state"
)

// StorageCluster represents config of storage cluster
type StorageCluster struct {
	Name   string       `json:"name"`
	Config state.Config `json:"config"`
}
