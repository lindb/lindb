package model

import (
	"encoding/json"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/planner/plan"
)

type RequestID string

type TaskID struct {
	RequestID RequestID `json:"requestId"`
	ID        int       `json:"id"`
}

type TaskRequest struct {
	TaskID     TaskID          `json:"taskId"`
	Partitions []int           `json:"partitions"`
	Fragment   json.RawMessage `json:"fragment,omitempty"`
}

type TaskResultSet struct {
	Page   *spi.Page       `json:"page,omitempty"`
	TaskID TaskID          `json:"taskId"`
	NodeID plan.PlanNodeID `json:"nodeId"`
	NoMore bool            `json:"noMore"`
}
