package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

var DB *pgxpool.Pool

func ConnectDB() {
	var err error
	databaseUrl := os.Getenv("DATABASE_URL")
	log.Println("Database URL:", databaseUrl)
	DB, err = pgxpool.Connect(context.Background(), databaseUrl)

	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect to database: %v \n", err))
	}

	log.Println("Successfully connected to database")
	createTables()
}

func createTables() {
	createUserTable := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		
		CREATE TABLE IF NOT EXISTS "user" (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			phone_number VARCHAR(255) UNIQUE
		);`

	_, err := DB.Exec(context.Background(), createUserTable)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to create user table: %v \n", err))
	}

}
