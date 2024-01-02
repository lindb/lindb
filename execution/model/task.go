package model

import (
	"encoding/json"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/planner/plan"
)

type OperatorKey struct {
	TaskID TaskID
	NodeID plan.PlanNodeID
}

type TaskID struct {
	RequestID string `json:"requestId"`
	ID        int    `json:"id"`
}

type TaskRequest struct {
	TaskID     TaskID          `json:"taskId"`
	Partitions []int           `json:"partitions"`
	Fragment   json.RawMessage `json:"fragment,omitempty"`
}

type TaskResultSet struct {
	TaskID TaskID          `json:"taskId"`
	NodeID plan.PlanNodeID `json:"nodeId"`
	Page   *spi.Page       `json:"page,omitempty"`
	NoMore bool            `json:"noMore"`
}
