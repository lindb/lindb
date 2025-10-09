package trace

import (
	"context"
	"fmt"

	"github.com/lindb/common/pkg/encoding"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/storage"
	tracestore "github.com/lindb/lindb/storage/trace"
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
		engine:       s.engine,
		table:        table,
		partitionIDs: partitions,

		predicate: predicate,
	}
}

type sourceConnector struct {
	engine storage.Engine

	table spi.TableHandle

	partitionIDs []int
	partitions   []*Partition

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
	page := types.NewPage()
	// timeColumn := types.NewColumn()
	// page.AppendColumn(types.ColumnMetadata{DataType: types.DTInt, Name: "timestamp"}, timeColumn)
	msgColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTJSON, Name: "callstack"}, msgColumn)
	// fieldsColumn := types.NewColumn()
	// page.AppendColumn(types.ColumnMetadata{DataType: types.DTString, Name: "fields"}, fieldsColumn)

	expr, ok := sc.predicate.(*tree.ComparisonExpression)
	var traceID string
	if ok {

		evalCtx := expression.NewEvalContext(context.TODO())
		traceID, _ = expression.EvalString(evalCtx, expr.Right)
		fmt.Println(traceID)
	}
	// TODO: check err

	for _, partition := range sc.partitions {
		for _, segment := range partition.segments {
			logSegment := segment.(*tracestore.Segment)
			// logIDs := logSegment.GetLogIDs(ns)
			fmt.Printf("trace Segment=%v\n", logSegment)
			logData, err := logSegment.GetTrace(traceID)
			if err != nil {
				fmt.Printf("get log err:%v\n", err)
			} else if len(logData) > 0 {
				for _, msg := range logData {
					FilterTracesByTraceID(traceID, msg, msgColumn)
				}
			}
		}
	}
	// batchs := jaeger.ProtoFromTraces(out)
	// json := string(encoding.JSONMarshal(batchs))
	// msgColumn.AppendString(json)
	// req := ptraceotlp.NewExportRequestFromTraces(out)
	// json, err := req.MarshalJSON()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Println(string(json))
	// msgColumn.AppendString(string(json))

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
		db:        db,
		timeRange: logTable.GetTimeRange(),
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
							tableScan: tableScan,
							shard:     shard,
							segments:  segments,
						})
					}
				}
			}
		}
	}
	return
}

func FilterTracesByTraceID(traceID string, msg []byte, column *types.Column) {
	req := ptraceotlp.NewExportRequest()
	if err := req.UnmarshalProto(msg); err != nil {
		fmt.Println(err)
		return
	}
	traces := req.Traces()
	resourceSpans := traces.ResourceSpans()

	if resourceSpans.Len() == 0 {
		return
	}

	for i := 0; i < resourceSpans.Len(); i++ {
		rs := resourceSpans.At(i)
		callStack := TranslateResourceSpans(rs, traceID)
		if callStack != nil {
			column.AppendJSON(encoding.JSONMarshal(callStack))
		}
	}
}
