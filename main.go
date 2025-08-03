package main

import (
    "context"
    "database/sql"
    "log"
    "net"
    "path/filepath"

    _ "github.com/mattn/go-sqlite3"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    pb "github.com/hoang19020302/go-user-service/userpb"
    "github.com/hoang19020302/go-user-service/internal/db"
)

type server struct {
    pb.UnimplementedUserServiceServer
    queries *db.Queries
}

func (s *server) GetUserById(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
    user, err := s.queries.GetUserByID(ctx, int64(req.GetId()))
    if err != nil {
        return nil, err
    }

    return &pb.UserResponse{
        Id:    int32(user.ID),
        Name:  user.Name,
        Email: user.Email,
    }, nil
}

func main() {
    absPath, err := filepath.Abs("db/db.sqlite")
    if err != nil {
        log.Fatal("Failed to get absolute path:", err)
    }

    conn, err := sql.Open("sqlite3", absPath)
    if err != nil {
        log.Fatal("Failed to open database:", err)
    }
    defer conn.Close()

    queries := db.New(conn)

    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterUserServiceServer(grpcServer, &server{queries: queries})

    // Bật reflection
    reflection.Register(grpcServer)

    log.Println("gRPC server listening on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

