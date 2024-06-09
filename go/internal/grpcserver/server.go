package grpcserver

import (
	"context"
	"go-django/internal/database"
	"go-django/internal/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	dbClient *database.DbClient
}

func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := req.GetUser()
	log.Println("update", user)
	s.dbClient.DB.Exec(
		`UPDATE users
		SET name = ?, email = ?
		WHERE id = ?`,
		user.Name,
		user.Email,
		user.GetId(),
	)
	return &pb.UpdateUserResponse{
		User: user,
	}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Add logic to delete user from database here
	log.Println("delete", req)
	s.dbClient.DB.Exec("DELETE FROM users WHERE id = ?", req.Id)
	// For example, let's pretend we deleted the user
	return &pb.DeleteUserResponse{Success: true}, nil
}

func StartGRPCServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	dbClient := database.InitDB()
	defer dbClient.DB.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &server{dbClient: dbClient})

	log.Printf("gRPC server listening on %s", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
