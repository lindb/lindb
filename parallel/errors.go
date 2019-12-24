package parallel

import "errors"

var errUnmarshalPlan = errors.New("unmarshal physical plan error")
var errUnmarshalQuery = errors.New("unmarshal query statement error")
var errUnmarshalSuggest = errors.New("unmarshal metadata suggest statement error")
var errWrongRequest = errors.New("not found task of current node from physical plan")
var errNoSendStream = errors.New("not found send stream")
var errTaskSend = errors.New("send task request error")
var errNoDatabase = errors.New("not found database")
