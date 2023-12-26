package database

import (
	"database/sql"
	"fmt"
)

func Init(db *sql.DB) error {
	// DROP TABLE IF EXISTS counter, counter_timestamp, world;
	sqlScript := `
		CREATE TABLE IF NOT EXISTS counter_timestamp (
			id BIGSERIAL PRIMARY KEY,
			timestamp TIMESTAMP UNIQUE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS world (
			id SERIAL PRIMARY KEY,
			name VARCHAR(64) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS counter (
			id BIGSERIAL PRIMARY KEY,
			count INT NOT NULL,
			world_id INT NOT NULL,
			counter_timestamp_id BIGINT NOT NULL,
			CONSTRAINT fk_counter_timestamp FOREIGN KEY (counter_timestamp_id) REFERENCES counter_timestamp(id),
			CONSTRAINT fk_world FOREIGN KEY (world_id) REFERENCES world(id)
		);
	`

	_, err := db.Exec(sqlScript)
	if err != nil {
		return err
	}

	fmt.Println("Tables created successfully")
	return nil
}
