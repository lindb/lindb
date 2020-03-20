mkdir -p rpc/pbmock/common
mkdir -p rpc/pbmock/storage

mockgen github.com/lindb/lindb/rpc/proto/common TaskServiceClient,TaskService_HandleClient,TaskServiceServer,TaskService_HandleServer > rpc/pbmock/common/common_mock.pb.go

mockgen github.com/lindb/lindb/rpc/proto/storage WriteServiceClient,WriteService_WriteClient,WriteServiceServer,WriteService_WriteServer  > rpc/pbmock/storage/storage_mock.pb.go
