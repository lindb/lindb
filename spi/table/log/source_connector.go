package log

import (
	"context"
	"fmt"
	"sort"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/proto/gen/v1/flatLogV1"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	logproto "github.com/lindb/lindb/proto/log"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/storage"
	"github.com/lindb/lindb/storage/log"
	logstore "github.com/lindb/lindb/storage/log"
)

type sourceConnectorProvider struct {
	engine storage.Engine
}

func NewSourceConnectorProvider(engine storage.Engine) spi.SourceConnectorProvider {
	return &sourceConnectorProvider{
		engine: engine,
	}
}

// CreateSourceConnector implements spi.SourceConnectorProvider.
func (s *sourceConnectorProvider) CreateSourceConnector(ctx context.Context,
	table spi.TableHandle, partitions []int, columnMapping map[string]string,
	predicate tree.Expression,
	outputColumns []types.ColumnMetadata, assignments []*spi.ColumnAssignment,
) spi.SourceConnector {
	return &sourceConnector{
		ctx:          ctx,
		engine:       s.engine,
		table:        table,
		partitionIDs: partitions,
		predicate:    predicate,
	}
}

type sourceConnector struct {
	ctx    context.Context
	engine storage.Engine

	table        spi.TableHandle
	partitionIDs []int

	partitions []*Partition

	predicate tree.Expression
}

// Run implements spi.SourceConnector.
func (sc *sourceConnector) Run(output chan<- *types.Page) {
	tableScan := sc.buildTableScan()
	if tableScan == nil {
		fmt.Println("table scan is nil")
		return
	}
	sc.partitions = sc.findPartitions(tableScan, sc.partitionIDs)
	if len(sc.partitions) == 0 {
		fmt.Printf("table partition is nil,ids=%v\n", sc.partitionIDs)
		return
	}
	indexDB := tableScan.db.IndexDatabase()
	ns, err := indexDB.GetNamespaceID([]byte("ns"))
	if err != nil {
		fmt.Printf("rr1=%v\n", err)
		panic(ns)
	}
	tableScan.nsID = ns

	if sc.predicate != nil {
		fieldLookup := NewFieldValuesLookupVisitor(sc.ctx, tableScan)
		fieldLookup.Visit(sc.ctx, sc.predicate)
	}

	page := types.NewPage()
	timeColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTInt, Name: "timestamp"}, timeColumn)
	msgColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTString, Name: "_msg"}, msgColumn)
	fieldsColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTJSON, Name: "fields"}, fieldsColumn)

	total := 0

	// sort partitions desc
	sort.Slice(sc.partitions, func(i, j int) bool {
		return sc.partitions[i].paritition.PartitionTime() > sc.partitions[j].paritition.PartitionTime()
	})

	for _, partition := range sc.partitions {
		// sort segments desc
		sort.Slice(partition.segments, func(i, j int) bool {
			return partition.segments[i].SegmentTimeRange().Start > partition.segments[j].SegmentTimeRange().Start
		})

		for _, segment := range partition.segments {
			logSegment := segment.(*logstore.Segment)
			logIDs := logSegment.FindLogIDsByTimeRange(tableScan.timeRange)
			if sc.predicate != nil {
				rowLookup := NewRowLookupVisitor(tableScan, logSegment, logIDs)
				logIDsObj := rowLookup.Visit(sc.ctx, sc.predicate)
				if logIDsByPredicate, ok := logIDsObj.(*roaring.Bitmap); ok {
					logIDs.And(logIDsByPredicate)
				}
			}
			if logIDs == nil || logIDs.IsEmpty() {
				continue
			}
			fmt.Printf("logSegment=%v,log ids:%v,%v\n", logSegment, logIDs)
			it := logIDs.ReverseIterator()
			for it.HasNext() {
				logID := it.Next()
				logData, err := logSegment.GetLog(logID)
				if err != nil {
					fmt.Printf("get log err:%v\n", err)
				} else {
					log := &flatLogV1.Log{}
					log.Init(logData, flatbuffers.GetUOffsetT(logData))
					timeColumn.AppendInt(log.Timestamp())
					// TODO:
					msgColumn.AppendString(string(log.Message()))
					fIt := logproto.NewFieldIterator(log)
					fMap := make(map[string]string)
					for fIt.HasNext() {
						fMap[string(fIt.NextName())] = string(fIt.NextValue())
					}
					fieldsColumn.AppendJSON(encoding.JSONMarshal(fMap))
					total++

					if total >= 1000 {
						// limit return
						goto END
					}
				}
			}
		}
	}
END:

	output <- page
}

func (sc *sourceConnector) buildTableScan() *TableScan {
	logTable, ok := sc.table.(*TableHandle)
	if !ok {
		panic(fmt.Sprintf("metric provider not support table handle<%T>", sc.table))
	}
	db, ok := sc.engine.GetDatabase2(logTable.Database)
	if !ok {
		panic(fmt.Errorf("%w: %s", constants.ErrDatabaseNotFound, logTable.Database))
	}

	return &TableScan{
		db:        db.(*log.Database),
		timeRange: logTable.GetTimeRange(),

		predicate: sc.predicate,
	}
}

func (sc *sourceConnector) findPartitions(tableScan *TableScan, partitionIDs []int) (partitions []*Partition) {
	for _, id := range partitionIDs {
		shard, ok := tableScan.db.GetShard(models.ShardID(id))
		if ok {
			pList := shard.GetPartitions(tableScan.timeRange)
			fmt.Printf("partitions=%v\n", pList)
			if len(pList) > 0 {
				for _, partition := range pList {
					segments := partition.GetSegments(tableScan.timeRange)
					fmt.Printf("segments=%v\n", segments)
					if len(segments) > 0 {
						partitions = append(partitions, &Partition{
							tableScan:  tableScan,
							shard:      shard,
							paritition: partition,
							segments:   segments,
						})
					}
				}
			}
		}
	}
	return
}
