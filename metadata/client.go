package main

import (
	"fmt"
	"time"

	pb "../proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	timestampFormat = time.StampNano // "Jan _2 15:04:05.000"
)

func main() {

	conn, err := grpc.Dial("127.0.0.1:50001", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := pb.NewGreeterClient(conn)
	//生成metadata数据
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	//发送metadata数据（客户端和服务端都可以发送metadata）
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "hello, world"})
	if err == nil {
		fmt.Printf("Reply is %s\n", resp.Message)
	} else {
		fmt.Printf("call server error:%s\n", err)
	}

}
