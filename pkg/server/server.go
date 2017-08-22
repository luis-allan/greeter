package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/soheilhy/cmux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/thingful/greeter/pkg/greeter"
	"github.com/thingful/greeter/pkg/version"
)

const (
	port = ":5555"
)

type server struct {
	logger log.Logger
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	s.logger.Log("handler", "SayHello", "name", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "version", version.Version, "ts", log.DefaultTimestampUTC)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Log("error", err, "port", port, "msg", "failed to listen")
		os.Exit(1)
	}

	logger.Log("msg", "starting server")

	// set up our request multiplexer
	m := cmux.New(lis)
	grpcLis := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpLis := m.Match(cmux.Any())

	// create grpc server
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &server{logger: log.With(logger, "protocol", "grpc")})

	// create http mux and server
	h := http.NewServeMux()
	h.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		logger.Log("protocol", "http", "path", "/status")
		fmt.Fprintf(w, "ok")
	})

	httpServer := &http.Server{
		Handler: h,
	}

	eps := make(chan error, 2)

	// start the listeners for each protocol
	go func() { eps <- grpcServer.Serve(grpcLis) }()
	go func() { eps <- httpServer.Serve(httpLis) }()

	logger.Log("port", port, "starting multiplexed server")
	err = m.Serve()

	var failed bool
	if err != nil {
		logger.Log("error", err, "msg", "cmux serve error")
		failed = true
	}

	var i int
	for err := range eps {
		if err != nil {
			logger.Log("error", err, "msg", "protocol serve error")
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
