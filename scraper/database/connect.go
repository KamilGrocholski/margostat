package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"scraper/config"
)

func Connect() (*sql.DB, error) {
	var db *sql.DB
	var err error

	for attempt := 1; attempt <= int(config.SCRAPER_DB_CONNECTION_MAX_ATTEMPTS); attempt++ {
		db, err = sql.Open("postgres", config.DATABASE_URL)
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
