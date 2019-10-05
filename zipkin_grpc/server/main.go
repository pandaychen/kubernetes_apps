package main

//grpc-server with zipkin

import (
	"log"
	"net"
	"time"

	//"github.com/opentracing/opentracing-go"

	//zipkin "github.com/openzipkin/zipkin-go-opentracing"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	global "../global"
	pb "../proto/add"
	pbcache "../proto/cache"
)

const (
	port    = ":50051"
	address = "localhost:50052"
)

// anthor rpc method

func GetCache(ctx context.Context, tracer opentracing.Tracer, id int32) int32 {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pbcache.NewCacheClient(conn)

	// Contact the server and print out its response.
	r, err := c.Get(ctx, &pbcache.CacheRequest{Id: id})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return -1
	}

	return r.GetResult()
}

// server is used to implement helloworld.GreeterServer.
type AddServer struct{}

// SayHello implements helloworld.GreeterServer
func (s *AddServer) DoAdd(ctx context.Context, in *pb.AddRequest) (*pb.AddReply, error) {

	log.Printf("input %d %d", in.GetNum1(), in.GetNum2())

	time.Sleep(time.Duration(10) * time.Millisecond)

	tracer := opentracing.GlobalTracer()

	//调用其他RPC方法（zipkin-tracer传递）
	val := GetCache(ctx, tracer, in.GetNum1())
	log.Printf("cache value %d", val)

	return &pb.AddReply{Result: val + in.GetNum2()}, nil
}

func main() {

	/*
		collector, err := zipkin.NewHTTPCollector("http://localhost:9411/api/v1/spans")
		if err != nil {
			log.Fatal(err)
			return
		}

		tracer, err := zipkin.NewTracer(
			zipkin.NewRecorder(collector, false, "localhost:0", "grpc_server"),
			zipkin.ClientServerSameSpan(true),
			zipkin.TraceID128Bit(true),
		)
		if err != nil {
			log.Fatal(err)
			return
		}
		opentracing.InitGlobalTracer(tracer)
	*/

	// set up a span reporter
	reporter := zipkinhttp.NewReporter(global.Zipkin_addr)
	defer reporter.Close()

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(global.Service_name, global.Service_endpoint)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)

	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads())))
	pb.RegisterAddServer(s, &AddServer{})
	s.Serve(lis)
}
