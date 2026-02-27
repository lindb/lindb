package traces

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow/memory"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/lindb/arrow/pkg/model"
	"github.com/lindb/arrow/pkg/traces"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"

	writerpkg "github.com/lindb/lindb/app/broker/write/writer"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
)

type writer struct {
	ctx          context.Context
	database     writerpkg.DatabaseAccessor[*model.Span]
	intervalCalc timeutil.IntervalCalculator
}

func NewWriter(
	ctx context.Context,
	databaseCfg models.Database,
	database writerpkg.DatabaseAccessor[*model.Span],
) *writer {
	return &writer{
		ctx:          ctx,
		intervalCalc: databaseCfg.Option.Intervals[0].Interval.Calculator(),
		database:     database,
	}
}

func (w *writer) Write(ctx context.Context, data []byte, encoding constants.EncodingType) error {
	shards := w.database.GetShards()
	numOfShards := uint64(len(shards))
	if numOfShards == 0 {
		return constants.ErrNoAvailableStorageNode
	}
	req := ptraceotlp.NewExportRequest()
	if err := req.UnmarshalProto(data); err != nil {
		return err
	}

	reqTraces := req.Traces()
	resourceSpans := reqTraces.ResourceSpans()

	if resourceSpans.Len() == 0 {
		return nil
	}

	for i := 0; i < resourceSpans.Len(); i++ {
		resource := resourceSpans.At(i)
		scopeSpans := resource.ScopeSpans()
		for j := 0; j < scopeSpans.Len(); j++ {
			scopeSpan := scopeSpans.At(j)
			spans := scopeSpan.Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				mSpan := &model.Span{
					Span:     span,
					Scope:    scopeSpan,
					Resource: resource,
				}

				shard := shards[mSpan.Hash()%numOfShards]
				// calculate segment time by span start timestamp
				segmentTime := w.intervalCalc.CalcFamilyTime(int64(span.StartTimestamp()) / 1000_000)
				segment := shard.GetOrCreateSegment(segmentTime, func() larrow.EntryBuilder[*model.Span] {
					return traces.NewSpanBuilder(memory.NewGoAllocator())
				})

				// write span to segment
				segment.Write(ctx, mSpan)
			}
		}
	}
	return nil
}
