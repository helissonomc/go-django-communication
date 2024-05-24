package main

import (
    "go-django/internal/database"
    "go-django/internal/routers"
	"log"
    "net/http"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
    database.InitDB()
    defer database.DB.Close()
    router := routers.InitRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
