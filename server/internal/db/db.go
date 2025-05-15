package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

var DB *pgxpool.Pool

func InitDB() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dsn := "postgres://" + dbUser + ":" + dbPassword + "@localhost:5432/" + dbName + "?sslmode=disable"
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DB, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	err = DB.Ping(ctx)
	if err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	fmt.Println("Connected to PostgreSQL successfully.")
}
