package postgres

import (
	"database/sql"
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

func (pa *PostgresAdapter) Query(query string, args ...any) (*sql.Rows, error) {
	return pa.DB.Query(query, args...)
}

func (pa *PostgresAdapter) Close() error {
	return pa.DB.Close()
}

// For local dev/testing.
func (pa *PostgresAdapter) CheckPost() {
	// Query data
	rows, err := pa.DB.Query("SELECT id, name, email, age FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over rows
	for rows.Next() {
		var id int
		var name, email string
		var age int

		err := rows.Scan(&id, &name, &email, &age)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("User: ID=%d, Name=%s, Email=%s, Age=%d\n", id, name, email, age)
	}

	// Check for errors after iteration
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}
