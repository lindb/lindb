package rpc

import "github.com/eleme/lindb/rpc/pkg/common"

const (
	OK int32 = iota
	ERR
)

func BuildResponse(code int32, msg string, data []byte) *common.Response {
	return &common.Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func ResponseOK() *common.Response {
	return BuildResponse(OK, "", nil)
}

func ResponseOKWithData(data []byte) *common.Response {
	return BuildResponse(OK, "", data)
}

func ResponseError(msg string) *common.Response {
	return BuildResponse(ERR, msg, nil)
}
