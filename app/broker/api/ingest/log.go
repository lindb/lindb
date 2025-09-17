package ingest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lindb/common/log"
	"github.com/lindb/common/models"
	"github.com/lindb/common/pkg/http"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
)

const Path = "/logs"

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
	route.POST(Path, w.Write)
	route.PUT(Path, w.Write)
}

func (w *Log) Write(c *gin.Context) {
	if err := w.write(c); err != nil {
		http.Error(c, err)
	}
}

func (w *Log) write(c *gin.Context) error {
	var logs []models.Log
	fmt.Println("write log...")
	if err := c.ShouldBind(&logs); err != nil {
		return err
	}
	fmt.Println(logs)
	rb := log.CreateRowBuilder()
	for _, l := range logs {
		rb.AddMessage([]byte(l.Message)).
			AddTimestamp(l.Timestamp)
		for k, v := range l.Fields {
			rb.AddField([]byte(k), []byte(v))
		}
		data, _ := rb.Build()
		fmt.Println(string(data))
		if err := w.deps.CM.WriteMsg(c.Request.Context(), "log_test", data); err != nil {
			return err
		}
		rb.Reset()
	}

	// contentType := strings.ToLower(strings.Trim(c.Request.Header.Get(headers.ContentType), " "))
	// var err error
	// var rows *metric.BrokerBatchRows
	// var writeType string
	// switch {
	// case strings.HasPrefix(contentType, constants.ContentTypeFlat):
	// 	rows, err = flat.Parse(c.Request, enrichedTags, param.Namespace, limits)
	// 	writeType = "flat"
	// case strings.HasPrefix(contentType, constants.ContentTypeJson):
	// 	rows, err = influx.Parse(c.Request, enrichedTags, param.Namespace, limits)
	// 	writeType = "json"
	// case strings.HasPrefix(contentType, constants.ContentTypeProto):
	// 	writeType = "proto"
	// 	rows, err = proto.Parse(c.Request, enrichedTags, param.Namespace, limits)
	// default:
	// 	err = fmt.Errorf("not support content type: %s, only support %s/%s/%s", contentType,
	// 		constants.ContentTypeFlat, constants.ContentTypeProto, constants.ContentTypeJson)
	// }
	// if err != nil {
	// 	return err
	// }
	return nil
}
