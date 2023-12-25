package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	maxAttempts, parseErr := strconv.ParseInt(os.Getenv("SCRAPER_DB_CONNECTION_MAX_ATTEMPTS"), 10, 32)
	if parseErr != nil {
		return nil, parseErr
	}
	dbURL := os.Getenv("DATABASE_URL")
	if len(dbURL) == 0 {
		panic("dbURL length is 0")
	}

	var db *sql.DB
	var err error

	for attempt := 1; attempt <= int(maxAttempts); attempt++ {
		db, err = sql.Open("postgres", dbURL)
		if err == nil {
			break
		}

		fmt.Printf("Error connecting to the database (attempt %d): %v\n", attempt, err)
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	if err != nil {
		fmt.Println("Failed to connect to the database after multiple attempts.")
		return nil, err
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Error pinging the database:", err)
		return nil, err
	}

	fmt.Println("Connected to db")
	return db, nil
}
