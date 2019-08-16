package rpc

import pb "github.com/lindb/lindb/rpc/proto/common"

//go:generate mockgen -source ./task_receiver.go -destination=./task_receiver_mock.go -package=rpc

// TaskReceiver represents the task result receiver
type TaskReceiver interface {
	// Receive receives the task result
	Receive(req *pb.TaskResponse) error
}
