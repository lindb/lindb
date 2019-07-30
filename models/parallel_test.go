package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhysicalPlan(t *testing.T) {
	physicalPlan := NewPhysicalPlan(Root{Indicator: "1.1.1.3:8000", NumOfTask: 1})
	physicalPlan.AddLeaf(Leaf{
		BaseNode: BaseNode{
			Parent:    "1.1.1.2:8000",
			Indicator: "1.1.1.1:9000",
		},
		Receivers: []Node{{IP: "1.1.1.5", Port: 8000}},
		ShardIDs:  []int32{1, 2, 4},
	})
	physicalPlan.AddIntermediate(Intermediate{
		BaseNode: BaseNode{
			Parent:    "1.1.1.3:8000",
			Indicator: "1.1.1.2:8000",
		},
		NumOfTask: 1,
	})
	physicalPlan.Database = "test_db"

	assert.Equal(t, PhysicalPlan{
		Database: "test_db",
		Root:     Root{Indicator: "1.1.1.3:8000", NumOfTask: 1},
		Intermediates: []Intermediate{{
			BaseNode: BaseNode{
				Parent:    "1.1.1.3:8000",
				Indicator: "1.1.1.2:8000",
			},
			NumOfTask: 1}},
		Leafs: []Leaf{{
			BaseNode: BaseNode{
				Parent:    "1.1.1.2:8000",
				Indicator: "1.1.1.1:9000",
			},
			Receivers: []Node{{IP: "1.1.1.5", Port: 8000}},
			ShardIDs:  []int32{1, 2, 4},
		}},
	}, *physicalPlan)
}
