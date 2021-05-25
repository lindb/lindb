// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package query

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
)

func TestBrokerPlan_Wrong_Case(t *testing.T) {
	plan := newBrokerPlan("sql", models.Database{}, nil, models.Node{}, nil)
	// storage nodes cannot be empty
	err := plan.Plan()
	assert.Equal(t, errNoAvailableStorageNode, err)

	storageNodes := map[string][]int32{"1.1.1.1:8000": {1, 2, 4}}
	// wrong sql
	plan = newBrokerPlan("sql", models.Database{}, storageNodes, models.Node{}, nil)
	err = plan.Plan()
	assert.NotNil(t, err)
}

func TestBrokerPlan_wrong_database_interval(t *testing.T) {
	storageNodes := map[string][]int32{"1.1.1.1:9000": {1, 2, 4}, "1.1.1.2:9000": {3, 5, 6}}
	currentNode := generateBrokerActiveNode("1.1.1.3", 8000)
	// no group sql
	plan := newBrokerPlan("select f from cpu",
		models.Database{Option: option.DatabaseOption{Interval: "s"}},
		storageNodes, currentNode.Node, nil)
	err := plan.Plan()
	assert.Error(t, err)
}

func TestBrokerPlan_No_GroupBy(t *testing.T) {
	storageNodes := map[string][]int32{"1.1.1.1:9000": {1, 2, 4}, "1.1.1.2:9000": {3, 5, 6}}
	currentNode := generateBrokerActiveNode("1.1.1.3", 8000)
	// no group sql
	plan := newBrokerPlan("select f from cpu",
		models.Database{Option: option.DatabaseOption{Interval: "10s"}},
		storageNodes, currentNode.Node, nil)
	err := plan.Plan()
	assert.NoError(t, err)

	p := plan.(*brokerPlan)
	assert.Equal(t, 0, len(p.intermediateNodes))
	physicalPlan := models.NewPhysicalPlan(models.Root{Indicator: "1.1.1.3:8000", NumOfTask: 2})
	physicalPlan.AddLeaf(models.Leaf{
		BaseNode: models.BaseNode{
			Parent:    "1.1.1.3:8000",
			Indicator: "1.1.1.1:9000",
		},
		Receivers: []models.Node{currentNode.Node},
		ShardIDs:  []int32{1, 2, 4},
	})
	physicalPlan.AddLeaf(models.Leaf{
		BaseNode: models.BaseNode{
			Parent:    "1.1.1.3:8000",
			Indicator: "1.1.1.2:9000",
		},
		Receivers: []models.Node{currentNode.Node},
		ShardIDs:  []int32{3, 5, 6},
	})
	assert.Equal(t, physicalPlan.Root, p.physicalPlan.Root)
	assert.Equal(t, 2, len(p.physicalPlan.Leafs))
	assert.Equal(t, 0, len(p.physicalPlan.Intermediates))
}

func TestBrokerPlan_GroupBy(t *testing.T) {
	storageNodes := map[string][]int32{
		"1.1.1.1:9000": {1, 2, 4},
		"1.1.1.2:9000": {3, 6, 9},
		"1.1.1.3:9000": {5, 7, 8},
		"1.1.1.4:9000": {10, 13, 15},
		"1.1.1.5:9000": {11, 12, 14},
	}
	currentNode := generateBrokerActiveNode("1.1.1.3", 8000)
	plan := newBrokerPlan(
		"select f from cpu group by host",
		models.Database{Option: option.DatabaseOption{Interval: "10s"}},
		storageNodes,
		currentNode.Node,
		[]models.ActiveNode{
			generateBrokerActiveNode("1.1.1.1", 8000),
			generateBrokerActiveNode("1.1.1.2", 8000),
			currentNode,
			generateBrokerActiveNode("1.1.1.4", 8000),
		})
	err := plan.Plan()
	if err != nil {
		t.Fatal(err)
	}
	p := plan.(*brokerPlan)
	assert.Equal(t, 3, len(p.intermediateNodes))
	physicalPlan := p.physicalPlan
	assert.Equal(t, models.Root{Indicator: "1.1.1.3:8000", NumOfTask: 3}, physicalPlan.Root)
	assert.Equal(t, 3, len(physicalPlan.Intermediates))
	for _, intermediate := range physicalPlan.Intermediates {
		assert.Equal(t, "1.1.1.3:8000", intermediate.Parent)
		assert.Equal(t, int32(5), intermediate.NumOfTask)
	}
	assert.Equal(t, 5, len(physicalPlan.Leafs))
	storageNodes2 := make(map[string][]int32)
	for _, leaf := range physicalPlan.Leafs {
		storageNodes2[leaf.Indicator] = leaf.ShardIDs
		assert.Equal(t, 3, len(leaf.Receivers))
	}
	assert.Equal(t, storageNodes, storageNodes2)
}

func TestBrokerPlan_GroupBy_Less_StorageNodes(t *testing.T) {
	storageNodes := map[string][]int32{
		"1.1.1.1:9000": {1, 2, 4},
		"1.1.1.2:9000": {3, 5, 6},
	}
	currentNode := generateBrokerActiveNode("1.1.1.3", 8000)
	plan := newBrokerPlan(
		"select f from cpu group by host",
		models.Database{Option: option.DatabaseOption{Interval: "10s"}},
		storageNodes,
		currentNode.Node,
		[]models.ActiveNode{
			generateBrokerActiveNode("1.1.1.1", 8000),
			generateBrokerActiveNode("1.1.1.2", 8000),
			currentNode,
			generateBrokerActiveNode("1.1.1.4", 8000),
		})
	err := plan.Plan()
	if err != nil {
		t.Fatal(err)
	}
	p := plan.(*brokerPlan)
	assert.Equal(t, 3, len(p.intermediateNodes))
	physicalPlan := p.physicalPlan
	assert.Equal(t, models.Root{Indicator: "1.1.1.3:8000", NumOfTask: 3}, physicalPlan.Root)
	assert.Equal(t, 3, len(physicalPlan.Intermediates))
	for _, intermediate := range physicalPlan.Intermediates {
		assert.Equal(t, "1.1.1.3:8000", intermediate.Parent)
		assert.Equal(t, int32(2), intermediate.NumOfTask)
	}
	assert.Equal(t, 2, len(physicalPlan.Leafs))
	storageNodes2 := make(map[string][]int32)
	for _, leaf := range physicalPlan.Leafs {
		storageNodes2[leaf.Indicator] = leaf.ShardIDs
		assert.Equal(t, 3, len(leaf.Receivers))
	}
	assert.Equal(t, storageNodes, storageNodes2)
}

func TestBrokerPlan_GroupBy_Same_Broker(t *testing.T) {
	storageNodes := map[string][]int32{"1.1.1.1:9000": {1, 2, 4}}
	currentNode := generateBrokerActiveNode("1.1.1.3", 8000)

	// current node = active node
	plan := newBrokerPlan(
		"select f from cpu group by host",
		models.Database{Option: option.DatabaseOption{Interval: "10s"}},
		storageNodes,
		currentNode.Node,
		[]models.ActiveNode{currentNode})
	err := plan.Plan()
	if err != nil {
		t.Fatal(err)
	}
	p := plan.(*brokerPlan)
	assert.Equal(t, 0, len(p.intermediateNodes))
	physicalPlan := models.NewPhysicalPlan(models.Root{Indicator: "1.1.1.3:8000", NumOfTask: 1})
	physicalPlan.AddLeaf(models.Leaf{
		BaseNode: models.BaseNode{
			Parent:    "1.1.1.3:8000",
			Indicator: "1.1.1.1:9000",
		},
		Receivers: []models.Node{currentNode.Node},
		ShardIDs:  []int32{1, 2, 4},
	})
	assert.Equal(t, physicalPlan, p.physicalPlan)
}

func TestBrokerPlan_GroupBy_No_Broker(t *testing.T) {
	storageNodes := map[string][]int32{"1.1.1.1:9000": {1, 2, 4}}
	currentNode := generateBrokerActiveNode("1.1.1.3", 8000)

	// only one storage node
	plan := newBrokerPlan(
		"select f from cpu group by host",
		models.Database{Option: option.DatabaseOption{Interval: "10s"}},
		storageNodes,
		currentNode.Node,
		nil)
	err := plan.Plan()
	if err != nil {
		t.Fatal(err)
	}
	p := plan.(*brokerPlan)
	assert.Equal(t, 0, len(p.intermediateNodes))
	physicalPlan := models.NewPhysicalPlan(models.Root{Indicator: "1.1.1.3:8000", NumOfTask: 1})
	physicalPlan.AddLeaf(models.Leaf{
		BaseNode: models.BaseNode{
			Parent:    "1.1.1.3:8000",
			Indicator: "1.1.1.1:9000",
		},
		Receivers: []models.Node{currentNode.Node},
		ShardIDs:  []int32{1, 2, 4},
	})
	assert.Equal(t, physicalPlan, p.physicalPlan)
}

func TestBrokerPlan_GroupBy_One_StorageNode(t *testing.T) {
	storageNodes := map[string][]int32{"1.1.1.1:9000": {1, 2, 4}}
	currentNode := generateBrokerActiveNode("1.1.1.3", 8000)

	// only one storage node
	plan := newBrokerPlan(
		"select f from cpu group by host",
		models.Database{Option: option.DatabaseOption{Interval: "10s"}},
		storageNodes,
		currentNode.Node,
		[]models.ActiveNode{
			generateBrokerActiveNode("1.1.1.1", 8000),
			generateBrokerActiveNode("1.1.1.2", 8100),
			currentNode,
			generateBrokerActiveNode("1.1.1.4", 8200),
		})
	err := plan.Plan()
	if err != nil {
		t.Fatal(err)
	}
	p := plan.(*brokerPlan)
	assert.Equal(t, 0, len(p.intermediateNodes))
	physicalPlan := models.NewPhysicalPlan(models.Root{Indicator: "1.1.1.3:8000", NumOfTask: 1})
	physicalPlan.AddLeaf(models.Leaf{
		BaseNode: models.BaseNode{
			Parent:    "1.1.1.3:8000",
			Indicator: "1.1.1.1:9000",
		},
		Receivers: []models.Node{currentNode.Node},
		ShardIDs:  []int32{1, 2, 4},
	})
	assert.Equal(t, physicalPlan, p.physicalPlan)
}

func generateBrokerActiveNode(ip string, port int) models.ActiveNode {
	return models.ActiveNode{Node: models.Node{IP: ip, Port: uint16(port)}}
}
