package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/app"
	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
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
	db := database.InitDB(mutex)
	app := app.App{
		DB:       db,
		Session:  session.InitSession(),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		WG:       &sync.WaitGroup{},
		MU:       mutex,
		Models:   database.New(db),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: routes.InitRoutes(&app),
	}

	go listenForShutdown(&app)
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

func listenForShutdown(app *app.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	shutDown(app)
	os.Exit(0)
}

func shutDown(app *app.App) {
	app.InfoLog.Println("Cleanup before shutdown, wait to background routines")
	app.WG.Wait()
	if app.DB != nil {
		app.DB.Close()
	}
	app.InfoLog.Println("Application shutdown")
}
