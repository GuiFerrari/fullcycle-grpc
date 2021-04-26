package main

import (
	"net"
	"log"

	"github.com/guiferrari/fullcycle-grpc/pb"
	"github.com/guiferrari/fullcycle-grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")

	if err != nil {
		log.Fatalf("Cound not connect: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, services.NewUserService())
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Cound not server: %v", err)
	}
}