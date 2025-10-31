package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/silent-observer/go-tickets/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedCommentServiceServer
}

func (s *server) Create(_ context.Context, in *pb.Comment) (*pb.CommentFullId, error) {
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
	pb.RegisterCommentServiceServer(s, &server{})
	fmt.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}