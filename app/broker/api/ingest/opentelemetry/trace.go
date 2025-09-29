package opentelemetry

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
	"github.com/lindb/common/pkg/http"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
)

const TracePath = "/opentelemetry/traces"

type Trace struct {
	deps *depspkg.HTTPDeps

	statistics struct {
		flat   *linmetric.BoundHistogram
		proto  *linmetric.BoundHistogram
		influx *linmetric.BoundHistogram
	}
}

func NewTrace(deps *depspkg.HTTPDeps) *Trace {
	ingestStatistics := metrics.NewCommonIngestionStatistics()

	return &Trace{
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
func (w *Trace) Register(route gin.IRoutes) {
	route.POST(TracePath, w.Write)
	route.PUT(TracePath, w.Write)
}

func (w *Trace) Write(c *gin.Context) {
	if err := w.write(c); err != nil {
		fmt.Println(err)
		http.Error(c, err)
	}
}

func (w *Trace) write(c *gin.Context) error {
	fmt.Println(c.Request.Header.Get(headers.ContentType))
	fmt.Printf("write trace...%v==\n", c)

	// 读取请求体（OTLP trace数据是 protobuf 格式）
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	defer c.Request.Body.Close()

	// 反序列化 protobuf 数据到 TraceServiceRequest 结构体
	req := ptraceotlp.NewExportRequest()
	if err := req.UnmarshalProto(body); err != nil {
		return err
	}

	// 简单打印接收到的 Trace Request 信息
	fmt.Printf("Received trace request: %+v\n", req)
	traces := req.Traces()
	spans := traces.ResourceSpans()
	for i := range spans.Len() {
		s := spans.At(i)
		resource := s.Resource()
		fmt.Printf("attributes: %+v\n", resource.Attributes())
	}
	if err := w.deps.CM.WriteMsg(c.Request.Context(), "trace_test", body); err != nil {
		return err
	}

	return nil
}
