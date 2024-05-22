package main

import (
	"database/sql"
	"fmt"
	"go-django/middleware"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	// Get database connection info from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Printf("Connected to the database successfully!")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})
	// Adding methods to a router
	r.HandleFunc("/books/{title}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World! book"))
	}).Methods("POST")

	// Restrict the request handler to specific hostnames or subdomains.
	r.HandleFunc("/p/books/{title}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World! site"))
	}).Host("www.mybookstore.com")

	// Restrict the request handler to http/https.
	// r.HandleFunc("/secure", SecureHandler).Schemes("https")
	// r.HandleFunc("/insecure", InsecureHandler).Schemes("http")
	bookrouter := r.PathPrefix("/prefix").Subrouter()
	bookrouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World! site"))
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}
