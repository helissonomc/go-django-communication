package grpcclient

import (
	"context"
	"go-django/internal/pb"
	"log"
	"time"

	"google.golang.org/grpc"
)

type client struct {
	conn          *grpc.ClientConn
	serviceClient pb.UserServiceClient
}

type GrpClientInterface interface {
	CreateUser(context.Context, *pb.User) (*pb.CreateUserResponse, error)
	UpdateUser(ctx context.Context, user *pb.User) (*pb.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, id int32) (*pb.DeleteUserResponse, error)
	Close()
}

func (c *client) CreateUser(ctx context.Context, user *pb.User) (*pb.CreateUserResponse, error) {
	log.Printf("Client is: %+v", c.serviceClient)
	return c.serviceClient.CreateUser(ctx, &pb.CreateUserRequest{User: user})
}

func (c *client) UpdateUser(ctx context.Context, user *pb.User) (*pb.UpdateUserResponse, error) {
	return c.serviceClient.UpdateUser(ctx, &pb.UpdateUserRequest{User: user})
}

func (c *client) DeleteUser(ctx context.Context, id int32) (*pb.DeleteUserResponse, error) {
	return c.serviceClient.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
}

func (client *client) Close() {
	client.conn.Close()
}

func NewClient(address string) (GrpClientInterface, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	serviceClient := pb.NewUserServiceClient(conn)
	log.Println("gRPC client initialized", serviceClient)

	return &client{
		conn:          conn,
		serviceClient: serviceClient,
	}, nil
}
