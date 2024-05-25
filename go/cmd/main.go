package main

import (
	"go-django/internal/database"
	"go-django/internal/grpc_client"
	"go-django/internal/routers"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	database.InitDB()
	defer database.DB.Close()
	router := routers.InitRouter()
	grpc_client.InitGRPCClient()

	log.Printf("Listening to localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
