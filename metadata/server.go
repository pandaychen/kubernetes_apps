package main

import (
	pb "../proto"
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var host = "127.0.0.1"

var (
	ServiceName = flag.String("ServiceName", "hello_service", "service name")
	Port        = flag.Int("Port", 50001, "listening port")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *Port))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	} else {
		fmt.Printf("listen at:%d\n", *Port)
	}
	defer lis.Close()

	s := grpc.NewServer()
	defer s.GracefulStop()

	pb.RegisterGreeterServer(s, &server{})
	addr := fmt.Sprintf("%s:%d", host, *Port)
	fmt.Printf("server addr:%s\n", addr)

	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %s", err)
	}
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// 服务端尝试接收metadata数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Printf("get metadata error")
	}
	if t, ok := md["timestamp"]; ok {
		fmt.Printf("timestamp from metadata:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	}
	//fmt.Printf("%v: Receive is %s\n", time.Now(), in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
