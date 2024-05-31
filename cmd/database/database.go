package database

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const MaxConnections = 10

func InitDB(mu *sync.Mutex) *sql.DB {
	conn := connectDB(mu)
	if conn == nil {
		panic("cannot connect to DB")
	}

	return conn
}

func connectDB(mu *sync.Mutex) *sql.DB {
	counts := 0

	dsn := os.Getenv("DSN")

	mu.Lock()
	defer func() {
		mu.Unlock()
	}()
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Reconnecting again to db...")
			counts++
		} else {
			connection.SetMaxOpenConns(MaxConnections)
			return connection
		}

		if counts > 5 {
			return nil
		}

		time.Sleep(time.Duration(counts) * time.Second) //Exponential backoff
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("db connection open failed with error: %s\n", err.Error())

		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("db ping failed with error: %s\n", err.Error())

		return nil, err
	}

	return db, nil
}
