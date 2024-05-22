package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-django/middleware"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

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
    r.HandleFunc("/books/{title}", func(w http.ResponseWriter, r *http.Request){
        w.Write([]byte("Hello, World! book"))
    }).Methods("POST")
    
    // Restrict the request handler to specific hostnames or subdomains.
    r.HandleFunc("/p/books/{title}", func(w http.ResponseWriter, r *http.Request){
        w.Write([]byte("Hello, World! site"))
    }).Host("www.mybookstore.com")
    
    // Restrict the request handler to http/https.
    // r.HandleFunc("/secure", SecureHandler).Schemes("https")
    // r.HandleFunc("/insecure", InsecureHandler).Schemes("http")
    bookrouter := r.PathPrefix("/books").Subrouter()
    bookrouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
        w.Write([]byte("Hello, World! site"))
    })
	log.Fatal(http.ListenAndServe(":8080", r))
}
