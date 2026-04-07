package main

import (
	"log/slog"
	"net"
	"os"

	"github.com/wutipong/albums/worker/service"
	"github.com/wutipong/albums/worker/service/pb"
	"google.golang.org/grpc"
)

//go:generate protoc --go_out=. --go-grpc_out=. -I/workspaces/grpc worker.proto

func main() {
	address := os.Getenv("WORKER_ADDRESS")
	if address == "" {
		slog.Error("invalid address")
		return
	}
	lis, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("unable to start server", slog.String("error", err.Error()))
		return
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterWorkerServiceServer(grpcServer, &service.WorkerServiceServer{})
	grpcServer.Serve(lis)
}
