package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDatabase() (*sql.DB, error) {
	var db *sql.DB
	var err error

	/*env := os.Getenv("ENVIRONMENT")

	if env != "production" {
		database := os.Getenv("DB_DATABASE")
		password := os.Getenv("DB_PASSWORD")
		username := os.Getenv("DB_USERNAME")
		port := os.Getenv("DB_PORT")
		host := os.Getenv("DB_HOST")

		localDSN := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, username, password, database,
		)

		db, err = sql.Open("pgx", localDSN)
		if err != nil {
			log.Fatalf("Failed to open local database: %v", err)
		}
	} else {*/
	connStr := os.Getenv("CONNECTION_STRING")
	if connStr == "" {
		log.Fatal("CONNECTION_STRING must be set into production environment.")
	}

	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to open production database: %v", err)
	}
	//}

	return db, nil
}
