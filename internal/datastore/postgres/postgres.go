package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func CheckPost() {
	var (
		host     = os.Getenv("POSTGRES_HOST")
		port     = os.Getenv("POSTGRES_HOST_PORT")
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname   = os.Getenv("POSTGRES_DB")
	)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	log.Printf("Connecting to: host=%s port=%s user=%s dbname=%s", host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	// Query data
	rows, err := db.Query("SELECT id, name, email, age FROM users")
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
