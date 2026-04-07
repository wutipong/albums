package main

import (
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/wutipong/albums/worker/service"
	"google.golang.org/grpc"

	pb "github.com/wutipong/albums/worker/service/definition"
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
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterWorkerServiceServer(grpcServer, &service.WorkerServiceServer{})
	grpcServer.Serve(lis)
}
