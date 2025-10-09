package opentelemetry

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/lindb/common/log"
	"github.com/lindb/common/pkg/http"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
)

const LogPath = "/opentelemetry/logs"

type Log struct {
	deps *depspkg.HTTPDeps

	statistics struct {
		flat   *linmetric.BoundHistogram
		proto  *linmetric.BoundHistogram
		influx *linmetric.BoundHistogram
	}
}

func NewLog(deps *depspkg.HTTPDeps) *Log {
	ingestStatistics := metrics.NewCommonIngestionStatistics()

	return &Log{
		deps: deps,
		statistics: struct {
			flat   *linmetric.BoundHistogram
			proto  *linmetric.BoundHistogram
			influx *linmetric.BoundHistogram
		}{
			flat:   ingestStatistics.Duration.WithTagValues("flat"),
			proto:  ingestStatistics.Duration.WithTagValues("proto"),
			influx: ingestStatistics.Duration.WithTagValues("influx"),
		},
	}
}

// Register adds the log ingest url route.
func (w *Log) Register(route gin.IRoutes) {
	route.POST(LogPath, w.Write)
	route.PUT(LogPath, w.Write)
}

func (w *Log) Write(c *gin.Context) {
	if err := w.write(c); err != nil {
		fmt.Println(err)
		http.Error(c, err)
	}
}

func (w *Log) write(c *gin.Context) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	defer c.Request.Body.Close()

	req := plogotlp.NewExportRequest()
	if err := req.UnmarshalProto(body); err != nil {
		return err
	}
	rb := log.CreateRowBuilder()

	logs := req.Logs()

	rLogs := logs.ResourceLogs()
	for i := range rLogs.Len() {
		log := rLogs.At(i)
		attr := log.Resource().Attributes()
		scopeLogs := log.ScopeLogs()
		for j := range scopeLogs.Len() {
			sl := scopeLogs.At(j)
			lrs := sl.LogRecords()
			for k := range lrs.Len() {
				lr := lrs.At(k)
				rb.AddMessage([]byte(lr.Body().AsString())).
					AddTimestamp(lr.Timestamp().AsTime().UnixMilli())
				rb.AddField([]byte("level"), []byte(lr.SeverityText()))
				attr.Range(func(k string, v pcommon.Value) bool {
					if v.AsString() != "" {
						rb.AddField([]byte(k), []byte(v.AsString()))
					}
					return true
				})
				lr.Attributes().Range(func(k string, v pcommon.Value) bool {
					if v.AsString() != "" {
						rb.AddField([]byte(k), []byte(v.AsString()))
					}
					return true
				})

				data, _ := rb.Build()
				dData := make([]byte, len(data))
				copy(dData, data)
				if err := w.deps.CM.WriteMsg(c.Request.Context(), "log_test", dData); err != nil {
					return err
				}
				rb.Reset()
			}
		}
	}

	return nil
}
