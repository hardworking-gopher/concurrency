package main

import (
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
	"time"
)

func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Fatal("failed to connect to db")
	}

	return conn
}

func connectToDB() *sql.DB {
	connStr := os.Getenv(envPgConnStr)

	for i := 0; i < pgConnRetries; i++ {
		db, err := openDB(connStr)
		if err != nil {
			log.Println("failed to open pgx conn", err.Error())
			log.Println("retrying...")

			time.Sleep(pgConnRetryDelay)
		} else {
			return db
		}
	}

	return nil
}

func openDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
