package handler

import (
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc/proto/field"
)

// normal case
func TestTcpHandler_HandleConn_Normal(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	cm := replication.NewMockChannelManager(ctl)

	h := NewTCPHandler(cm)

	in, out := net.Pipe()

	done := make(chan struct{})

	go func() {
		if err := h.Handle(out); err != nil {
			t.Error(err)
		}
		done <- struct{}{}
	}()

	// first metric list
	metricList := buildMetricList(1)
	metricListBytes, err := metricList.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	writer := stream.NewBufferWriter(nil)
	writer.PutInt32(int32(len(metricListBytes)))
	writer.PutBytes(metricListBytes)

	cm.EXPECT().Write(gomock.Any(), metricList).Return(nil)

	bytes, err := writer.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	n, err := in.Write(bytes)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(bytes) {
		t.Fatal("should not happen")
	}

	// second metric list
	metricList2 := buildMetricList(1)
	metricListBytes2, err := metricList2.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	writer.Reset()
	writer.PutInt32(int32(len(metricListBytes2)))
	writer.PutBytes(metricListBytes2)

	cm.EXPECT().Write(gomock.Any(), metricList2).Return(nil)

	bytes2, err := writer.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	n2, err := in.Write(bytes2)
	if err != nil {
		t.Fatal(err)
	}

	if n2 != len(bytes) {
		t.Fatal("should not happen")
	}

	if err := in.Close(); err != nil {
		t.Fatal(err)
	}

	<-done
}

// length bytes不足4字节
func TestTcpHandler_SizeBytesNotEnough(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	cm := replication.NewMockChannelManager(ctl)

	h := NewTCPHandler(cm)

	in, out := net.Pipe()

	go func() {
		if err := h.Handle(out); err != nil {
			t.Error(err)
		}
	}()

	writer := stream.NewBufferWriter(nil)
	writer.PutByte(1)

	bytes, err := writer.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	n, err := in.Write(bytes)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(bytes) {
		t.Fatal("should not happen")
	}

	if err := in.Close(); err != nil {
		t.Fatal(err)
	}
}

// data bytes不足 length bytes表示的长度
func TestTcpHandler_DataBytesNotEnough(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	cm := replication.NewMockChannelManager(ctl)

	h := NewTCPHandler(cm)

	in, out := net.Pipe()

	go func() {
		if err := h.Handle(out); err != nil {
			t.Error(err)
		}
	}()

	writer := stream.NewBufferWriter(nil)
	writer.PutInt32(2)
	writer.PutByte(1)

	bytes, err := writer.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	n, err := in.Write(bytes)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(bytes) {
		t.Fatal("should not happen")
	}

	if err := in.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestTcpHandler_UnmarshalFail(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	cm := replication.NewMockChannelManager(ctl)

	h := NewTCPHandler(cm)

	in, out := net.Pipe()

	done := make(chan struct{})

	go func() {
		err := h.Handle(out)
		if err == nil {
			t.Error("should be error")
		}
		logger.GetLogger("broker", "TCPHandler").Error("handler error", logger.Error(err))
		done <- struct{}{}
	}()

	// first metric list
	metricList := buildMetricList(1)
	metricListBytes, err := metricList.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// fake bit error
	metricListBytes[0]++

	writer := stream.NewBufferWriter(nil)
	writer.PutInt32(int32(len(metricListBytes)))
	writer.PutBytes(metricListBytes)

	bytes, err := writer.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	n, err := in.Write(bytes)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(bytes) {
		t.Fatal("should not happen")
	}

	<-done
}

func buildMetricList(value float64) *field.MetricList {
	return &field.MetricList{
		Metrics: []*field.Metric{{
			Name:      "name",
			Timestamp: time.Now().Unix() * 1000,
			Tags:      map[string]string{"tagKey": "tagVal"},
			Fields: []*field.Field{{
				Name: "sum",
				Field: &field.Field_Sum{
					Sum: &field.Sum{
						Value: value,
					},
				},
			}},
		}}}
}
