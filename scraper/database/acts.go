package database

import (
	"database/sql"
	"time"

	"scraper/utils"
)

type CounterInsert struct {
	Count int
	World string
}

type CountResult struct {
	Timestamp string
	Count     int
}

func GetAllWorldNames(db *sql.DB) ([]string, error) {
	query := `
		SELECT name 
		FROM world
		ORDER BY name
	`
	result, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	out := []string{}
	for result.Next() {
		var name string
		err := result.Scan(&name)
		if err != nil {
			return nil, err
		}
		out = append(out, name)
	}

	return out, nil
}

func GetGlobalTimeline(db *sql.DB) ([]CountResult, error) {
	query := `
		SELECT
			ct.timestamp,
			SUM(c.count) AS count
		FROM
			counter c
		INNER JOIN
			counter_timestamp ct ON c.counter_timestamp_id = ct.id
		GROUP BY
			ct.timestamp
		ORDER BY
			ct.timestamp;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []CountResult{}

	for rows.Next() {
		var countResult CountResult
		var time time.Time

		err := rows.Scan(
			&time,
			&countResult.Count,
		)
		if err != nil {
			return nil, err
		}

		countResult.Timestamp = utils.FormatDateToJavascriptISOString(time)

		out = append(out, countResult)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

// func GetWorldTimeline(db *sql.DB, world string, timeStart time.Time, timeEnd time.Time) ([]CountResult, error) {
func GetWorldTimeline(db *sql.DB, world string) ([]CountResult, error) {
	query := `
		SELECT
			c.count,
			ct.timestamp
		FROM
			counter c
		INNER JOIN
			world w ON c.world_id = w.id
		INNER JOIN
			counter_timestamp ct ON c.counter_timestamp_id = ct.id
		WHERE
			w.name = $1
	`

	rows, err := db.Query(query, world)
	if err != nil {
		return nil, err
	}

	out := []CountResult{}

	for rows.Next() {
		countResult := CountResult{}
		var time time.Time
		err := rows.Scan(&countResult.Count, &time)
		if err != nil {
			return nil, err
		}
		countResult.Timestamp = utils.FormatDateToJavascriptISOString(time)
		out = append(out, countResult)
	}

	return out, nil
}

func InsertMultipleCounters(db *sql.DB, data []CounterInsert) error {
	currentTimestamp := time.Now()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO counter (count, world_id, counter_timestamp_id)
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	timestampID, err := getOrCreateTimestampID(tx, currentTimestamp)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, counter := range data {
		worldID, err := getOrCreateWorldID(tx, counter.World)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = stmt.Exec(counter.Count, worldID, timestampID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func getOrCreateTimestampID(tx *sql.Tx, timestamp time.Time) (int64, error) {
	var timestampID int64
	err := tx.QueryRow("SELECT id FROM counter_timestamp WHERE timestamp = $1", timestamp).Scan(&timestampID)

	if err == sql.ErrNoRows {
		err = tx.QueryRow("INSERT INTO counter_timestamp (timestamp) VALUES ($1) RETURNING id", timestamp).Scan(&timestampID)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	return timestampID, nil
}

func getOrCreateWorldID(tx *sql.Tx, worldName string) (int64, error) {
	var worldID int64
	err := tx.QueryRow("SELECT id FROM world WHERE name = $1", worldName).Scan(&worldID)

	if err == sql.ErrNoRows {
		err = tx.QueryRow("INSERT INTO world (name) VALUES ($1) RETURNING id", worldName).Scan(&worldID)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	return worldID, nil
}
