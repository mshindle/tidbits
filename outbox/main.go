package outbox

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-faker/faker/v4"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/mshindle/tidbits/outbox/internal/repository/sqlite"
	"gitlab.com/mshindle/tidbits/outbox/internal/user"
)

func Execute(ctx context.Context) error {
	// connect to the database
	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// create the repos & service
	store := sqlite.NewStore(db)
	err = store.InitializeDB()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	userService := user.NewService(store, store, store, realClock{})

	// generate the dummy data
	var person Person
	_ = faker.FakeData(&person)
	_, err = userService.Create(ctx, person.Name, person.Email)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func connectDB() (*sql.DB, error) {
	// --- Database Setup (for demonstration, replace with your actual connection) ---
	db, err := sql.Open("sqlite3", "file:outbox.sqlite")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("connected to database")
	return db, nil
}

type Person struct {
	Name  string `faker:"name"`
	Email string `faker:"email"`
}

type realClock struct{}

func (rc realClock) Now() time.Time {
	return time.Now()
}
