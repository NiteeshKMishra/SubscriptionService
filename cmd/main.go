package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/emailer"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/routes"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/session"
)

const PORT = 8080
const SECRET_FILE = "secrets.env"

func main() {
	log.Println("Welcome to subscription service")

	err := readSecrets()
	if err != nil {
		log.Panicf("error in parsing secrets %s", err.Error())
	}

	mutex := &sync.Mutex{}
	wt := &sync.WaitGroup{}
	db := database.InitDB(mutex)
	app := app.App{
		DB:       db,
		Session:  session.InitSession(),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		WG:       wt,
		MU:       mutex,
		Models:   database.New(db),
		Mailer:   emailer.NewMailer(wt),
	}

	go app.ListenForMail()
	go app.ListenForShutdown()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: routes.InitRoutes(&app),
	}

	app.InfoLog.Printf("Starting web server on port %d\n", PORT)
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func readSecrets() error {
	requiredSecrets := []string{
		"DB_DSN",
		"REDIS_HOST",
		"PLANS",
		"ADMIN_EMAIL",
		"ADMIN_PASSWORD",
		"MAIL_DOMAIN",
		"MAIL_HOST",
		"MAIL_PORT",
		"MAIL_ENCRYPTION",
		"MAIL_FROM",
		"MAIL_ADDRESS",
	}
	file, err := os.Stat(SECRET_FILE)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("os.stat failed with Error %s", err.Error())
		return err
	}

	if err == nil && file.Size() > 0 {
		err := godotenv.Load(SECRET_FILE)
		if err != nil {
			log.Panicf("unable to parse secrets with error %s", err.Error())
		}
	}

	for _, secret := range requiredSecrets {
		if os.Getenv(secret) == "" {
			log.Panicf("required secret %s not present", secret)
		}
	}

	return nil
}
