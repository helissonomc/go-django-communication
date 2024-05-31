package main

import (
	"go-django/internal/grpcserver"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	address := ":50051"
	grpcserver.StartGRPCServer(address)
}
