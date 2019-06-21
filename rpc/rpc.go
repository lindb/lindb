package rpc

import (
	"sync"

	"go.uber.org/zap"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/rpc/pkg/batch"
	"github.com/eleme/lindb/rpc/pkg/broker"
	"github.com/eleme/lindb/rpc/pkg/common"
)

var (
	// lock to guard register toBatchFuncMap, fromBatchFuncMap
	lock             = sync.Mutex{}
	toBatchFuncMap   = map[common.RequestType]ToBatchFunc{}
	fromBatchFuncMap = map[common.RequestType]FromBatchFunc{}
)

// Request includes a specific Type and a corresponding rpc request pointer
type Request struct {
	Type              common.RequestType
	WritePointRequest *broker.WritePointsRequest
}

// ToBatchFunc convert request to batch request
type ToBatchFunc func(req *Request) *batch.BatchRequest_Request

// RegisterToBatchRequestFunc registers a requestFunc for an requestType
// Nil requestFunc and registering the same requestType more than once will result in a fatal
func RegisterToBatchRequestFunc(requestType common.RequestType, requestFunc ToBatchFunc) {
	if requestFunc == nil {
		logger.GetLogger().Fatal("RequestFunc is nil", zap.Stack("stack"))
	}

	lock.Lock()
	defer lock.Unlock()
	_, ok := toBatchFuncMap[requestType]
	if ok {
		logger.GetLogger().Fatal("ToBatchRequest already registered",
			zap.String("requestType", requestType.String()),
			zap.Stack("stack"))
	}
	toBatchFuncMap[requestType] = requestFunc
}

// ToBatchRequest converts Request to BatchRequest according to registered converted func,
// unknown request type will result in a fatal
func (req *Request) ToBatchRequest() *batch.BatchRequest_Request {
	toBatchReqFunc, ok := toBatchFuncMap[req.Type]
	if !ok {
		logger.GetLogger().Fatal("ToBatchFunc not registered",
			zap.String("requestType", req.Type.String()),
			zap.Stack("stack"))
	}
	return toBatchReqFunc(req)
}

// Response includes a specific RequestType and a corresponding response pointer
type Response struct {
	Type                common.RequestType
	WritePointsResponse *broker.WritePointsResponse
}

// FromBatchFunc converts BatchResponse to Response
type FromBatchFunc func(res *batch.BatchResponse_Response) *Response

// RegisterFromBatchFunc registers a responseFunc to converts batchResponse to response for a specific requestType
// Nil responseFunc and register a requestType for more than once will result in a fatal
func RegisterFromBatchFunc(requestType common.RequestType, responseFunc FromBatchFunc) {
	if responseFunc == nil {
		logger.GetLogger().Fatal("ResponseFunc in nil", zap.Stack("stack"))
	}
	lock.Lock()
	defer lock.Unlock()
	_, ok := fromBatchFuncMap[requestType]
	if ok {
		logger.GetLogger().Fatal("FromBatchResponse already registered",
			zap.String("requestType", requestType.String()),
			zap.Stack("stack"))
	}
	fromBatchFuncMap[requestType] = responseFunc
}

// FromBatchResponse converts batchResponse to response
func FromBatchResponse(res *batch.BatchResponse_Response) *Response {
	responseFunc, ok := fromBatchFuncMap[res.RequestType]
	if !ok {
		logger.GetLogger().Fatal("ToBatchFunc not registered",
			zap.String("requestType", res.RequestType.String()),
			zap.Stack("stack"))
	}
	return responseFunc(res)
}

// BuildRequestContext helps to build a RequestContext
func BuildRequestContext(requestType common.RequestType) *common.RequestContext {
	return &common.RequestContext{
		Type: requestType,
	}
}

// BuildResponseContext helps to build a ResponseContext
func BuildResponseContext(msg string) *common.ResponseContext {
	return &common.ResponseContext{
		Msg: msg,
	}
}
