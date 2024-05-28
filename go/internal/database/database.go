package database

import (
	"database/sql"
	"go-django/internal/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
)


type DbClient struct {
    DB *sql.DB
}

func NewDbClient(db *sql.DB) *DbClient{
    return &DbClient{DB: db}        
}

func (db *DbClient) createUsersTable() {
	query := `CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	);`

	_, err := db.DB.Exec(query)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}
}

func InitDB() *DbClient {
	cfg := config.GetConfig()

	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
    log.Printf("Connected to the database successfully!")
    dbClient := NewDbClient(db)
    dbClient.createUsersTable()
    return dbClient
}
