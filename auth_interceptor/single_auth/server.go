package main

import (
	pb "../../proto"
	"flag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"fmt"
)

var (
	port = flag.String("p", ":8972", "port")
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("call SayHello")
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//set Interceptor
	//s := grpc.NewServer(grpc.StreamInterceptor(StreamServerInterceptor),
	//	grpc.UnaryInterceptor(UnaryServerInterceptor1))
	s := grpc.NewServer(grpc.UnaryInterceptor(composeUnaryServerInterceptors(UnaryServerInterceptor1, UnaryServerInterceptor2)))
	pb.RegisterGreeterServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//compose two inceptors
// multi interceptors
func composeUnaryServerInterceptors(interceptor, next grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	fmt.Println("call composeUnaryServerInterceptors")
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		return interceptor(ctx, req, info,
			func(nextCtx context.Context, nextReq interface{}) (interface{}, error) {

				return next(nextCtx, nextReq, info, handler)
			})
	}
}

func UnaryServerInterceptor1(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("before handling interceptor1. Info: %+v", info)
	fmt.Println("interceptor1:",req,handler)
	resp, err := handler(ctx, req)
	log.Printf("after handling interceptor1. resp: %+v", resp)

	return resp, err
}

func UnaryServerInterceptor2(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("before handling interceptor2. Info: %+v", info)
	fmt.Println("interceptor2:",req,handler)
	resp, err := handler(ctx, req)
	log.Printf("after handling interceptor2. resp: %+v", resp)
	return resp, err
}

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("before handling. Info: %+v", info)
	err := handler(srv, ss)
	log.Printf("after handling. err: %v", err)
	return err
}
