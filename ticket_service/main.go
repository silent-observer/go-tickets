package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/silent-observer/go-tickets/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedTicketServiceServer
}

func (s *server) Create(_ context.Context, in *pb.Ticket) (*pb.TicketFullId, error) {
	log.Printf("Received %s", in.String())
	return in.GetId(), nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(s, healthcheck)
	pb.RegisterTicketServiceServer(s, &server{})
	fmt.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}