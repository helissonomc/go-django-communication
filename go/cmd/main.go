package main

import (
	"go-django/internal/controllers"
	"go-django/internal/database"
	"go-django/internal/grpcclient"
	"go-django/internal/routers"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)
  
func main() {
    dbClient := database.InitDB()
	defer dbClient.DB.Close() 
	grpcClient, err := grpcclient.NewClient("django-grpc-server:50052", "token_test")
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	log.Println("gRPC client initialized", grpcClient)

	userController := controllers.NewUserController(grpcClient, dbClient)
	router := routers.InitRouter(userController)

	log.Printf("Listening to localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
