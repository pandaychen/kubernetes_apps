package main

//GRPC-SERVER2

import (
	"log"
	"net"
	"time"

	"github.com/opentracing/opentracing-go"

	global "../global"
	pb "../proto/cache"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

// server is used to implement helloworld.GreeterServer.
type CacheServer struct{}

// SayHello implements helloworld.GreeterServer
// WHY no tracer???
func (s *CacheServer) Get(ctx context.Context, in *pb.CacheRequest) (*pb.CacheReply, error) {

	log.Printf("input %d", in.GetId())
	time.Sleep(time.Duration(50) * time.Millisecond)

	return &pb.CacheReply{Result: in.GetId() * 2}, nil
}

func main() {
	/*
		collector, err := zipkin.NewHTTPCollector("http://localhost:9411/api/v1/spans")
		if err != nil {
			log.Fatal(err)
			return
		}

		tracer, err := zipkin.NewTracer(
			zipkin.NewRecorder(collector, false, "localhost:0", "grpc_cache"),
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
	pb.RegisterCacheServer(s, &CacheServer{})
	s.Serve(lis)
}
