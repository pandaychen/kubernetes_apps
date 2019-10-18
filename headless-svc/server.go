package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "../proto"
)

type SearchService struct{}

func (s *SearchService) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("Recv:", r.Name)
	return &pb.HelloReply{Message: r.Name + " Server"}, nil
}

const PORT = "50051"

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &SearchService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}
