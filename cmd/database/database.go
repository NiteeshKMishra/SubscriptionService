package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path"
	"sync"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	"golang.org/x/crypto/bcrypt"
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

	dsn := os.Getenv("DB_DSN")

	for {
		connection, err := openDB(dsn, mu)
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

func openDB(dsn string, mu *sync.Mutex) (*sql.DB, error) {
	mu.Lock()
	defer func() {
		mu.Unlock()
	}()
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

	// Run migrations when db is available
	err = runMigrations(db)
	if err != nil {
		log.Printf("migrations failed with error: %s\n", err.Error())

		db.Close()

		return nil, err
	}

	// Populate data when db is available
	err = populateDB(db)
	if err != nil {
		log.Printf("populating db with error: %s\n", err.Error())

		db.Close()

		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	dbMigrations := &migrate.FileMigrationSource{Dir: path.Join("cmd", "migrations")}
	_, err := dbMigrations.FindMigrations()
	if err != nil {
		log.Printf("failed to find migrations with error: %s\n", err.Error())
		return err
	}

	out, err := migrate.Exec(db, "postgres", dbMigrations, migrate.Up)
	if err != nil {
		log.Printf("failed to run migrations with error: %s\n", err.Error())
		return err
	}
	log.Printf("added %d migrations to db successfully\n", out)

	return nil
}

func populateDB(db *sql.DB) error {
	//Populate plans
	var planData []InitPlan
	plans := os.Getenv("PLANS")
	err := json.Unmarshal([]byte(plans), &planData)
	if err != nil {
		return err
	}

	for _, plan := range planData {
		_, err = db.Exec(`INSERT INTO plans (plan_name, plan_amount, created_at, updated_at) VALUES
			($1, $2, now(), now()) ON CONFLICT ON CONSTRAINT plans_plan_name_key
			DO UPDATE SET (plan_amount, updated_at) = ($2, now())`,
			plan.Name,
			plan.Amount,
		)
		if err != nil {
			return err
		}
	}

	//Populate Admin user
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPass := os.Getenv("ADMIN_PASS")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPass), 12)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO users (email, password, first_name, last_name, user_active, is_admin, created_at, updated_at)
	VALUES($1, $2, 'System', 'Admin', true, true, now(), now()) ON CONFLICT ON CONSTRAINT users_email_key
	DO UPDATE SET (password, updated_at) = ($2, now())`,
		adminEmail,
		hashedPassword,
	)
	if err != nil {
		return err
	}
	log.Println("successfully populated data in db")

	return nil
}
