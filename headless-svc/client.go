package main

import (
	pb "../proto"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"log"
	"time"
)

func main() {
	var port = flag.Int("port", 12345, "connect port")
	var srvname = flag.String("svc", "headsvc", "service name")
	flag.Parse()
	if port == nil {
		return
	}

	scheme := fmt.Sprintf("dns:///%s:%d", *srvname, *port)

	conn, err := grpc.Dial(scheme, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	for i := 0; i < 100; i++ {
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{
			Name: "gRPC",
		})
		if err != nil {
			log.Fatalf("client.Search err: %v", err)
			time.Sleep(1 * time.Second)
		}

		log.Printf("resp: %s", resp.Message)
		time.Sleep(1 * time.Second)
	}
}
