package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/travis-james/DBCache/internal/config"
)

type PostgresAdapter struct {
	DB *sql.DB
}

func NewPostgres(cc *config.Config) (PostgresAdapter, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cc.DatastoreDBHost, cc.DatastoreDBPort, cc.DatastoreDBUser, cc.DatastoreDBPw, cc.DatastoreDBName)
	log.Printf("Connecting to: %s", psqlInfo) // TODO: REMOVE THIS EVENTUALLY
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected!")

	return PostgresAdapter{
		DB: db,
	}, nil
}

// QueryRows is set up to take whatever query, and assuming success, return the result as a []byte.
func (pa *PostgresAdapter) QueryRows(query string, args ...any) ([]byte, error) {
	// Use sql.Rows to get column names
	rows, err := pa.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Read column names
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Scan needs pointers. So we take the address of the elements of values to out into valauePtrs. valuePtrs is read, and thus data is written into values.
	values := make([]any, len(cols))
	valuePtrs := make([]any, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Cursor starts before the first row. Call next to see if there are any returned rows.
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	// map of db column name to value
	result := make(map[string]any)
	for i, col := range cols {
		val := values[i]
		// Convert []byte to string for readability
		if b, ok := val.([]byte); ok {
			result[col] = string(b)
		} else {
			result[col] = val
		}
	}

	return json.Marshal(result)
}

func (pa *PostgresAdapter) Close() error {
	return pa.DB.Close()
}
