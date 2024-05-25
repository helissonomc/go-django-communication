package grpc_client

import (
	"context"
	"go-django/internal/pb"
	"log"
	"time"

	"google.golang.org/grpc"
)

var client pb.UserServiceClient

func InitGRPCClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "django-grpc-server:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client = pb.NewUserServiceClient(conn)
	log.Println("gRPC client initialized", client)
}

func CreateUser(ctx context.Context, user *pb.User) (*pb.CreateUserResponse, error) {
	log.Printf("Client is: %+v", client)
	return client.CreateUser(ctx, &pb.CreateUserRequest{User: user})
}

func UpdateUser(ctx context.Context, user *pb.User) (*pb.UpdateUserResponse, error) {
	return client.UpdateUser(ctx, &pb.UpdateUserRequest{User: user})
}

func DeleteUser(ctx context.Context, id int32) (*pb.DeleteUserResponse, error) {
	return client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
}
