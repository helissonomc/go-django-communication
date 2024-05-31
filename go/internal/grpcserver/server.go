package grpcserver

import (
	"context"
	"go-django/internal/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := req.GetUser()
	// Add logic to update user in database here
	log.Println("update", user)
	// For example, let's pretend we updated the user and return it
	return &pb.UpdateUserResponse{
		User: user,
	}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Add logic to delete user from database here
	log.Println(req)
	// For example, let's pretend we deleted the user
	return &pb.DeleteUserResponse{}, nil
}

func StartGRPCServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &server{})

	log.Printf("gRPC server listening on %s", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
