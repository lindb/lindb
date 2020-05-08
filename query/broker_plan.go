package query

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
)

// brokerPlan represents the broker execute plan
type brokerPlan struct {
	sql               string
	query             *stmt.Query
	storageNodes      map[string][]int32
	currentBrokerNode models.Node
	brokerNodes       []models.ActiveNode
	intermediateNodes []models.Node
	databaseCfg       models.Database

	physicalPlan *models.PhysicalPlan
}

// newBrokerPlan creates broker execute plan
func newBrokerPlan(sql string, databaseCfg models.Database, storageNodes map[string][]int32,
	currentBrokerNode models.Node, brokerNodes []models.ActiveNode) Plan {
	return &brokerPlan{
		sql:               sql,
		databaseCfg:       databaseCfg,
		storageNodes:      storageNodes,
		currentBrokerNode: currentBrokerNode,
		brokerNodes:       brokerNodes,
	}
}

// Plan plans broker level query execute plan, there are some scenarios as below:
// 1) parse sql => stmt
// 2) build parallel exec tree
//    a) no group by => only need leafs
//    b) one storage node => only need leafs
//    c) no other active broker node => node need leafs
//    d) need intermediate computing nodes
func (p *brokerPlan) Plan() error {
	lenOfStorageNodes := len(p.storageNodes)
	if lenOfStorageNodes == 0 {
		return errNoAvailableStorageNode
	}

	query, err := sql.Parse(p.sql)
	if err != nil {
		return err
	}
	// set query statement
	p.query = query

	if query.Interval <= 0 {
		var interval timeutil.Interval
		if err := interval.ValueOf(p.databaseCfg.Option.Interval); err != nil {
			return err
		}
		query.Interval = interval
	}
	intervalVal := int64(query.Interval)
	p.query.TimeRange.Start = timeutil.Truncate(p.query.TimeRange.Start, intervalVal)
	p.query.TimeRange.End = timeutil.Truncate(p.query.TimeRange.End, intervalVal)

	root := p.currentBrokerNode

	p.buildIntermediateNodes()

	lenOfIntermediateNodes := len(p.intermediateNodes)

	if lenOfIntermediateNodes > 0 {
		// create parallel exec task
		p.physicalPlan = models.NewPhysicalPlan(models.Root{
			Indicator: (&root).Indicator(),
			NumOfTask: int32(lenOfIntermediateNodes)})

		p.buildIntermediates()
	} else {
		receivers := []models.Node{root}
		// create parallel exec task
		p.physicalPlan = models.NewPhysicalPlan(models.Root{
			Indicator: (&root).Indicator(),
			NumOfTask: int32(lenOfStorageNodes)})
		p.buildLeafs((&root).Indicator(), p.getStorageNodeIDs(), receivers)
	}

	return nil
}

// buildIntermediateNodes builds intermediate nodes if need
func (p *brokerPlan) buildIntermediateNodes() {
	if len(p.query.GroupBy) == 0 {
		return
	}
	if len(p.brokerNodes) == 0 {
		return
	}
	if len(p.storageNodes) == 1 {
		return
	}

	for _, brokerNode := range p.brokerNodes {
		if brokerNode.Node != p.currentBrokerNode {
			p.intermediateNodes = append(p.intermediateNodes, brokerNode.Node)
		}
	}
}

// getStorageNodeIDs returns storage node ids
func (p *brokerPlan) getStorageNodeIDs() []string {
	var storageNodeIDs []string
	for nodeID := range p.storageNodes {
		storageNodeIDs = append(storageNodeIDs, nodeID)
	}
	return storageNodeIDs
}

// buildIntermediates builds the intermediates computing layer
func (p *brokerPlan) buildIntermediates() {
	lenOfIntermediateNodes := len(p.intermediateNodes)
	lenOfStorageNodes := len(p.storageNodes)
	// calc degree of parallelism
	parallel := lenOfStorageNodes / lenOfIntermediateNodes
	if lenOfStorageNodes%lenOfIntermediateNodes != 0 {
		parallel++
	}

	storageNodeIDs := p.getStorageNodeIDs()

	var pos, end, idx = 0, 0, 0
	for {
		if pos > lenOfStorageNodes {
			break
		}
		end += parallel

		if end > lenOfStorageNodes {
			end = lenOfStorageNodes
		}

		intermediateNodeID := (&p.intermediateNodes[idx]).Indicator()

		// add intermediate task into parallel exec tree
		p.physicalPlan.AddIntermediate(models.Intermediate{
			BaseNode: models.BaseNode{
				Parent:    (&p.currentBrokerNode).Indicator(),
				Indicator: intermediateNodeID,
			},
			NumOfTask: int32(lenOfStorageNodes),
		})
		// add leaf tasks into parallel exec tree
		p.buildLeafs(intermediateNodeID, storageNodeIDs[pos:end], p.intermediateNodes)

		pos += parallel
		idx++
	}
}

// buildLeafs builds the leaf computing nodes based parent, nodes and result receivers
func (p *brokerPlan) buildLeafs(parentID string, nodeIDs []string, receivers []models.Node) {
	for _, nodeID := range nodeIDs {
		p.physicalPlan.AddLeaf(models.Leaf{
			BaseNode: models.BaseNode{
				Parent:    parentID,
				Indicator: nodeID,
			},
			ShardIDs:  p.storageNodes[nodeID],
			Receivers: receivers,
		})
	}
}
