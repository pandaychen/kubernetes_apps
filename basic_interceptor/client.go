package main

/*
	2020/02/25 14:02:49 before invoker. method: /proto.Greeter/SayHello, request:name:"world"
invoker calls succ
2020/02/25 14:02:49 after invoker. reply: message:"Hello world"
2020/02/25 14:02:49 Greeting: Hello world

*/

import (
	pb "../proto"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

var (
	address = flag.String("addr", "localhost:8972", "address")
	name    = flag.String("n", "world", "name")
)

func main() {
	flag.Parse()
	// 连接服务器
	conn, err := grpc.Dial(*address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(UnaryClientInterceptor),
		grpc.WithStreamInterceptor(StreamClientInterceptor))
	if err != nil {
		log.Fatalf("faild to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("before invoker. method: %+v, request:%+v", method, req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		fmt.Println("invoker calls errors", err)
	} else {
		fmt.Println("invoker calls succ")
	}
	log.Printf("after invoker. reply: %+v", reply)
	return err
}
func StreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Printf("before invoker. method: %+v, StreamDesc:%+v", method, desc)
	clientStream, err := streamer(ctx, desc, cc, method, opts...)
	log.Printf("before invoker. method: %+v", method)
	return clientStream, err
}
