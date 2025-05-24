package engine

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	_ "github.com/marcboeker/go-duckdb"
)

type DB struct {
	conn *sql.DB
}

func NewDB() (*DB, error) {
	conn, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("Failed to open DuckDB connection: %w", err)
	}
	
	err = conn.Ping()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Failed to ping DuckDB: %w", err)
	}

	log.Println("Successfully connected to DuckDB (in-memory).")
	return &DB{conn: conn}, nil
}


func (d *DB) Close() error {
	if d.conn != nil {
		log.Println("Closing DuckDB connection.")
		return d.conn.Close()
	}
	return nil
}


func (d *DB) ExecuteQuery(query string) (string, error) {
	log.Printf("Executing query %s\n", query)

	rows, err := d.conn.Query(query)
	if err != nil {
		return "", fmt.Errorf("DuckDB query failed: %w", err)
	}
	defer rows.Close()

	var result string
	if rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			return "", fmt.Errorf("Failed to scan result: %w", err)
		} 
	} else {
			err := rows.Err()
			if err != nil {
				return "", fmt.Errorf("Error iterating rows: %w", err)
		}
		result = "No rows returned."
	}

	time.Sleep(50 * time.Millisecond)

	return fmt.Sprintf("DuckDB processing: %s. Result: %s", query, result), nil
} 
