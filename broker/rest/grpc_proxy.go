package rest

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// Create the rpc gateway proxy server handler for rpc interface.
// you need to register all the interfaces that need to generate rpc proxy.
func CreateRPCProxyServerMux(addr string) (*runtime.ServeMux, context.CancelFunc, error) {
	fmt.Printf("gRpc server addr is %s", addr)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	if ctx == nil {
		fmt.Println("not exits")
	}
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	// register server endpoint,
	//opts := []grpc.DialOption{grpc.WithInsecure()}
	//gRPCServerEndpoint := flag.String("grpc-server-endpoint", addr, "gRPC server endpoint")
	//register the service proxy
	//err := pb.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	return mux, cancel, nil
}
