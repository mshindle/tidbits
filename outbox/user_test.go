package outbox

import (
	"context"
	"database/sql"
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestOutbox(t *testing.T) {
	// --- Database Setup (for demonstration, replace with your actual connection) ---
	db, err := sql.Open("sqlite3", "file:outbox.sqlite?cache=shared")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database!") // Ensure the 'users' table exists for the example

	_, err = db.Exec(CreateUsersTableSQL)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}
	log.Println("Users table checked/created.") // Ensure the 'outbox_events' table exists
	_, err = db.Exec(CreateOutboxEventsTableSQL)
	if err != nil {
		log.Fatalf("Failed to create outbox_events table: %v", err)
	}
	log.Println("Outbox events table checked/created.") // --- Example Usage ---

	userService := NewUserService(db)
	_, err = userService.Create(context.Background(), "Alice Smith", "alice@example.com")
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}
	_, err = userService.Create(context.Background(), "Bob Johnson", "bob@example.com")
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}
}
