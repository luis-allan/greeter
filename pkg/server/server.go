package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/soheilhy/cmux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/thingful/greeter/pkg/greeter"
)

const (
	port = ":5555"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// set up our request multiplexer
	m := cmux.New(lis)
	grpcLis := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpLis := m.Match(cmux.Any())

	// create grpc server
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &server{})

	// create http mux and server
	h := http.NewServeMux()
	h.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})

	httpServer := &http.Server{
		Handler: h,
	}

	eps := make(chan error, 2)

	// start the listeners for each protocol
	go func() { eps <- grpcServer.Serve(grpcLis) }()
	go func() { eps <- httpServer.Serve(httpLis) }()

	log.Println("starting multiplexed server on: %s", port)
	err = m.Serve()

	var failed bool
	if err != nil {
		log.Println("cmux serve error: %v", err)
		failed = true
	}

	var i int
	for err := range eps {
		if err != nil {
			log.Printf("protocol serve error: %v", err)
			failed = true
		}
		i++
		if i == cap(eps) {
			close(eps)
			break
		}
		if failed {
			os.Exit(1)
		}
	}
}
