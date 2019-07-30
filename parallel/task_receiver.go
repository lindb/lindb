package parallel

import pb "github.com/lindb/lindb/rpc/proto/common"

// TaskReceiver represents the sub task result receiver
type TaskReceiver interface {
	// Receive receives the sub task result, them merge those results
	Receive(req *pb.TaskResponse) error
}
