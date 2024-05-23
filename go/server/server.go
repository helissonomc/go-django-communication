package main

import (
	"database/sql"
	"fmt"
	"go-django/middleware"
	"log"
	"net/http"
    "time"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func tableExists(db *sql.DB, tableName string) (bool, error) {
    query := fmt.Sprintf("SHOW TABLES LIKE '%s';", tableName)
    var name string
    err := db.QueryRow(query).Scan(&name)
    if err != nil && err != sql.ErrNoRows {
        return false, err
    }
    return name == tableName, nil
}

func main() {
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	// Get database connection info from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
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
    { // Create a new table
        exists, err := tableExists(db, "users")
        if err != nil {
            panic(err)
        }
        if !exists {
            query := `
                CREATE TABLE users (
                    id INT AUTO_INCREMENT,
                    username TEXT NOT NULL,
                    password TEXT NOT NULL,
                    created_at DATETIME,
                    PRIMARY KEY (id)
                );`

            if _, err := db.Exec(query); err != nil {
                panic(err)
            }
            fmt.Println("Table created successfully.")
        }else {
            fmt.Println("Table already exists.")
        }
    }
    { // Insert a new user
        username := "johndoe"
        password := "secret"
        createdAt := time.Now()

        result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
        if err != nil {
            log.Fatal(err)
        }

        id, err := result.LastInsertId()
        fmt.Println(id)
    }

    { // Query a single user
        var (
            id        int
            username  string
            password  string
            createdAt time.Time
        )

        query := "SELECT id, username, password, created_at FROM users WHERE id = ?"
        if err := db.QueryRow(query, 1).Scan(&id, &username, &password, &createdAt); err != nil {
            log.Fatal(err)
        }

        fmt.Println(id, username, password, createdAt)
    }
    { // Query all users
        type user struct {
            id        int
            username  string
            password  string
            createdAt time.Time
        }

        rows, err := db.Query(`SELECT id, username, password, created_at FROM users`)
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()

        var users []user
        for rows.Next() {
            var u user

            err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
            if err != nil {
                log.Fatal(err)
            }
            users = append(users, u)
        }
        if err := rows.Err(); err != nil {
            log.Fatal(err)
        }

        log.Printf("%v", users[0].createdAt)
    }
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
