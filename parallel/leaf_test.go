package parallel

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/service"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func TestLeafTask_Process_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	storageService := service.NewMockStorageService(ctrl)
	executorFactory := NewMockExecutorFactory(ctrl)

	currentNode := models.Node{IP: "1.1.1.3", Port: 8000}
	processor := newLeafTask(currentNode, storageService, executorFactory, taskServerFactory)
	// unmarshal error
	err := processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: nil})
	assert.Equal(t, errUnmarshalPlan, err)

	plan, _ := json.Marshal(&models.PhysicalPlan{
		Leafs: []models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.4:8000"}}},
	})
	// wrong request
	err = processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: plan})
	assert.Equal(t, errWrongRequest, err)

	plan, _ = json.Marshal(&models.PhysicalPlan{
		Database: "test_db",
		Leafs:    []models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
	})

	// unmarshal query err
	engine := tsdb.NewMockEngine(ctrl)
	storageService.EXPECT().GetEngine(gomock.Any()).Return(engine)
	err = processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: plan, Payload: []byte{1, 2, 3}})
	assert.Equal(t, errUnmarshalQuery, err)

	plan, _ = json.Marshal(&models.PhysicalPlan{
		Database: "test_db",
		Leafs:    []models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
	})
	query := stmt.Query{MetricName: "cpu"}
	data := encoding.JSONMarshal(&query)
	// db not exist
	storageService.EXPECT().GetEngine(gomock.Any()).Return(nil)
	err = processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: plan, Payload: data})
	assert.Equal(t, errNoDatabase, err)

	storageService.EXPECT().GetEngine(gomock.Any()).Return(engine).MaxTimes(2)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(nil)
	err = processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: plan, Payload: data})
	assert.Equal(t, errNoSendStream, err)

	// test executor fail
	storageService.EXPECT().GetEngine(gomock.Any()).Return(engine).AnyTimes()
	exec := NewMockExecutor(ctrl)
	exec.EXPECT().Execute().Return(nil)
	exec.EXPECT().Error().Return(fmt.Errorf("err"))
	executorFactory.EXPECT().NewStorageExecutor(context.TODO(), gomock.Any(), gomock.Any(), gomock.Any()).Return(exec)
	serverStream := pb.NewMockTaskService_HandleServer(ctrl)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(serverStream)
	serverStream.EXPECT().Send(gomock.Any()).Return(nil)
	err = processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: plan, Payload: data})
	if err != nil {
		t.Fatal(err)
	}
}

func TestLeafProcessor_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	storageService := service.NewMockStorageService(ctrl)
	executorFactory := NewMockExecutorFactory(ctrl)

	currentNode := models.Node{IP: "1.1.1.3", Port: 8000}
	processor := newLeafTask(currentNode, storageService, executorFactory, taskServerFactory)
	engine := tsdb.NewMockEngine(ctrl)
	plan, _ := json.Marshal(&models.PhysicalPlan{
		Database: "test_db",
		Leafs:    []models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
	})
	query := stmt.Query{MetricName: "cpu"}
	data := encoding.JSONMarshal(&query)

	storageService.EXPECT().GetEngine(gomock.Any()).Return(engine)

	serverStream := pb.NewMockTaskService_HandleServer(ctrl)
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(serverStream)
	serverStream.EXPECT().Send(gomock.Any()).Return(nil)
	exec := NewMockExecutor(ctrl)
	exec.EXPECT().Execute().Return(nil)
	exec.EXPECT().Error().Return(nil)
	executorFactory.EXPECT().NewStorageExecutor(context.TODO(), gomock.Any(), gomock.Any(), gomock.Any()).Return(exec)
	err := processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: plan, Payload: data})
	if err != nil {
		t.Fatal(err)
	}

	storageService.EXPECT().GetEngine(gomock.Any()).Return(engine).AnyTimes()
	taskServerFactory.EXPECT().GetStream(gomock.Any()).Return(serverStream).AnyTimes()
	serverStream.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
	exec.EXPECT().Error().Return(nil).AnyTimes()

	resultCh := make(chan *series.TimeSeriesEvent)
	it := series.NewMockGroupedIterator(ctrl)
	it.EXPECT().HasNext().Return(true)
	it.EXPECT().Tags().Return(nil)
	fieldIt := series.NewMockFieldIterator(ctrl)
	it.EXPECT().Next().Return(fieldIt)
	fieldIt.EXPECT().Bytes().Return([]byte{1, 2, 3}, nil)
	fieldIt.EXPECT().FieldMeta().Return(field.Meta{Name: "f1"})
	it.EXPECT().HasNext().Return(false)
	go func() {
		resultCh <- &series.TimeSeriesEvent{
			Series: it,
		}
		time.AfterFunc(500*time.Millisecond, func() {
			close(resultCh)
		})
	}()
	exec.EXPECT().Execute().Return(resultCh)
	executorFactory.EXPECT().NewStorageExecutor(context.TODO(), gomock.Any(), gomock.Any(), gomock.Any()).Return(exec)
	_ = processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: plan, Payload: data})

	resultCh = make(chan *series.TimeSeriesEvent)
	it = series.NewMockGroupedIterator(ctrl)
	it.EXPECT().HasNext().Return(true)
	it.EXPECT().Tags().Return(nil)
	fieldIt = series.NewMockFieldIterator(ctrl)
	it.EXPECT().Next().Return(fieldIt)
	fieldIt.EXPECT().Bytes().Return(nil, fmt.Errorf("err"))
	go func() {
		resultCh <- &series.TimeSeriesEvent{
			Series: it,
		}
		resultCh <- &series.TimeSeriesEvent{
			Series: it,
		}
		time.AfterFunc(500*time.Millisecond, func() {
			close(resultCh)
		})
	}()
	exec.EXPECT().Execute().Return(resultCh)
	executorFactory.EXPECT().NewStorageExecutor(context.TODO(), gomock.Any(), gomock.Any(), gomock.Any()).Return(exec)
	_ = processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: plan, Payload: data})
}
