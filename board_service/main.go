package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/silent-observer/go-tickets/pgdb"
	pb "github.com/silent-observer/go-tickets/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var (
	port          = flag.Int("port", 50052, "The server port")
	postgres_conn = flag.String("postgres_conn",
		"postgres://postgres:password@localhost:5432/postgres?sslmode=disable",
		"PostgreSQL connection string")
)

type server struct {
	pb.UnimplementedBoardServiceServer
	db *bun.DB
}

func (s *server) Create(cxt context.Context, in *pb.Board) (*pb.BoardFullId, error) {
	log.Printf("Received %s", in.String())

	name := in.GetName()
	if name == "" {
		st := status.New(codes.InvalidArgument, "No board name specified")
		return nil, st.Err()
	}

	project_name := in.GetId().GetProject().GetId()
	if project_name == "" {
		st := status.New(codes.InvalidArgument, "No project name for board")
		return nil, st.Err()
	}

	project := pgdb.Project{}
	if err := s.db.NewSelect().
		Model(&project).
		Column("id").
		Where("name = ?", project_name).
		Limit(1).
		Scan(cxt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			st := status.Newf(codes.NotFound, "No project named \"%s\"", project_name)
			return nil, st.Err()
		} else {
			st := status.New(codes.Internal, err.Error())
			return nil, st.Err()
		}
	}

	board := pgdb.Board{Name: name, ProjectId: project.Id}
	_, err := s.db.NewInsert().
		Model(&board).
		Column("name", "project_id").
		Returning("id").
		Exec(cxt)
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		return nil, st.Err()
	}

	id := &pb.BoardFullId{
		Board:   &pb.BoardId{Id: board.Id},
		Project: &pb.ProjectId{Id: project_name},
	}
	return id, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(*postgres_conn)))
	pgdb.SetMaxOpenConns(25)
	pgdb.SetMaxIdleConns(10)
	pgdb.SetConnMaxLifetime(5 * time.Minute)
	pgdb.SetConnMaxIdleTime(5 * time.Minute)

	db := bun.NewDB(pgdb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	s := grpc.NewServer()
	pb.RegisterBoardServiceServer(s, &server{db: db})
	fmt.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
